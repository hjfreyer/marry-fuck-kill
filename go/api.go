package gomfk

import (
	"appengine"
	_ "appengine/datastore"
	"appengine/memcache"
	"appengine/urlfetch"
	_ "appengine/user"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	_ "gomfk/json_api"
	"gomfk/parse_args"
	"net/http"
	"net/url"
	_ "reflect"
	_ "strconv"
	_ "strings"
	_ "code.google.com/p/goprotobuf/proto"
	_ "text/template"
)

const RETRY_COUNT = 3

func WrapHandler(handler func(c *Context)) func(
	w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		cxt := appengine.NewContext(r)
		user := UserIdFromContext(r)

		c := Context{user, cxt, w, r}
		handler(&c)
	}
}

type Context struct {
	userId UserId
	cxt    appengine.Context
	w      http.ResponseWriter
	r      *http.Request
}

func (c *Context) Error(code int, err error) {
	c.Errorf(code, err.Error())
}

func (c *Context) Errorf(code int, f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	c.cxt.Errorf(msg)
	http.Error(c.w, msg, code)
}

type makeRequest struct {
	A makeRequest_Entity `parseArg:"a"`
	B makeRequest_Entity `parseArg:"b"`
	C makeRequest_Entity `parseArg:"c"`
}

type makeRequest_Entity struct {
	Name  string            `parseArg:"Name,required"`
	Image makeRequest_Image `parseArg:"Image"`
}

type makeRequest_Image struct {
	Url string `parseArg:"Url,required"`
	// ContentType string `parseArg:"ContentType"`
	// SourceUrl   string `parseArg:"sourceUrl"`
	// Salt        int64  `parseArg:"salt"`
	// Hash        int64  `parseArg:"hash"`
}

type makeResponse struct {
	Id int64 `json:"id"`
}

func ApiMakeHandler(c *Context) {
	var request makeRequest
	err := parse_args.ParseArgs(c.r, &request)
	if err != nil {
		c.Error(400, err)
		return
	}

	FetchOrDie := func(url string) *FetchedImage {
		image, err := FetchImage(c.cxt, url)
		maybePanic(err)
		return image
	}

	imageA := FetchOrDie(request.A.Image.Url)
	imageB := FetchOrDie(request.B.Image.Url)
	imageC := FetchOrDie(request.C.Image.Url)

	triple := TripleCreation{
		A: EntityCreation{
			Name:  request.A.Name,
			Image: imageA,
		},
		B: EntityCreation{
			Name:  request.B.Name,
			Image: imageB,
		},
		C: EntityCreation{
			Name:  request.C.Name,
			Image: imageC,
		},
		Creator: c.userId,
	}

	db := NewAppengineDataAccessor(c.cxt)
	tripleId, err := db.MakeTriple(triple)
	if err != nil {
		panic(err)
	}

	response := makeResponse{int64(tripleId)}
	responseMsg, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	c.w.Write(responseMsg)
}

type voteRequest struct {
	TripleId TripleId `parseArg:"triple_id,required"`
	Vote     Vote     `parseArg:"vote"`
}

func (v *voteRequest) Validate() error {
	if !v.Vote.IsValid() {
		return errors.New("Invalid vote: " + string(v.Vote))
	}
	return nil
}

func ApiVoteHandler(c *Context) {
	var request voteRequest
	err := parse_args.ParseArgs(c.r, &request)
	if err != nil {
		c.Error(400, err)
		return
	}

	db := NewAppengineDataAccessor(c.cxt)
	err = db.UpdateVote(request.TripleId, c.userId, request.Vote)

	if err != nil {
		panic(err)
	}
}

type HttpError struct {
	StatusCode int
	Message string
}

type Fetcher interface {
	Fetch(query string) *http.Response
}

func ImageSearch(query string, logger Logger, fetcher Fetcher) (ImageSearchResponse, HttpError) {
	var imageSearchJson []byte
	for attempt := 0; attempt < NUM_RETRIES; attempt++ {
		resp := fetcher.Fetch(query)
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var err error
			imageSearchJson, err = ioutil.ReadAll(resp.Body)
			CheckOk(err)
		}
	}
}

func CachedHandler

func ApiImageSearchHandler(c *Context) {
	query := c.r.FormValue("query")
	c.cxt.Infof(query)

	if item, err := memcache.Get(c.cxt, query); err == memcache.ErrCacheMiss {
		c.cxt.Infof("Memcache: miss")
	} else if err != nil {
		panic(err)
	} else {
		c.cxt.Infof("Memcache: hit")
		c.w.Write(item.Value)
		return
	}

	const API_BASE = "https://www.googleapis.com/customsearch/v1"

	values := url.Values{}
	values.Set("cx", "017343173679326196998:omutomvh_wi")
	values.Set("filter", "1")
	values.Set("safe", "medium")
	values.Set("searchType", "image")
	values.Set("fields", "items(image(contextLink,thumbnailLink))")
	values.Set("pp", "1")
	values.Set("key", "AIzaSyDbjy0CKTMV5DoJR07ZYF5w-KL7Ey5lyGY")
	values.Set("q", query)
	values.Set("userIp", c.r.RemoteAddr)

	url := API_BASE + "?" + values.Encode()

	c.cxt.Infof("API Query: %s", url)

	fetcher := urlfetch.Client(c.cxt)
	response, err := fetcher.Get(url)
	if err != nil {
		panic(err)
	}

	var resultBuffer bytes.Buffer
	if _, err := resultBuffer.ReadFrom(response.Body); err != nil {
		panic(err)
	}

	var res struct {
		Items []struct {
			Image struct {
				ContextLink   string
				ThumbnailLink string
			}
		}
	}
	if json.Unmarshal(resultBuffer.Bytes(), &res); err != nil {
		// TODO: 400 in this case.
		panic(err)
	}

	var ourResult ImageSearchResponse
	for _, item := range res.Items {
		img := item.Image
		ourResult.Image = append(ourResult.Image, &ImageMetadata{
			ContextUrl: &img.ContextLink,
			Url:     &img.ThumbnailLink,
		})
	}

	ourResultJson, err := json.Marshal(ourResult)
	if err != nil {
		panic(err)
	}

	item := &memcache.Item{
		Key:   query,
		Value: ourResultJson,
	}

	if err := memcache.Add(c.cxt, item); err == memcache.ErrNotStored {
		c.cxt.Warningf("item with key %q already exists", query)
	} else if err != nil {
		panic(err)
	} else {
		c.cxt.Infof("Cached successfully: %q", query)
	}

	if _, err := c.w.Write(ourResultJson); err != nil {
		panic(err)
	}
}
