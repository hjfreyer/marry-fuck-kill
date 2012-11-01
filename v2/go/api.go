package gomfk

import (
"appengine/urlfetch"
	"bytes"
_	"text/template"
	_"reflect"
	_"strconv"
	"appengine"
	_"appengine/datastore"
	_ "appengine/user"
	"encoding/json"
	"errors"
	 "fmt"
	_"gomfk/json_api"
	"gomfk/parse_args"
	"net/http"
	"net/url"
	_ "strings"
	"appengine/memcache"
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
	cxt appengine.Context
	w http.ResponseWriter
	r *http.Request
}

func (c *Context) Error(code int, err error) {
	c.Errorf(code, err.Error())
}

func (c *Context) Errorf(code int, f string, args... interface{}) {
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
	Name  string  `parseArg:"Name,required"`
	Image makeRequest_Image  `parseArg:"Image"`
}

type makeRequest_Image struct {
	Url         string `parseArg:"Url,required"`
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

	imageA, err := FetchImage(c.cxt, request.A.Image.Url)
	if err != nil {
		panic(err)
	}
	imageB, err := FetchImage(c.cxt, request.B.Image.Url)
	if err != nil {
		panic(err)
	}
	imageC, err := FetchImage(c.cxt, request.C.Image.Url)
	if err != nil {
		panic(err)
	}

	triple := TripleCreation{
	A: EntityCreation {
		Name: request.A.Name,
		Image: imageA,
		},
	B: EntityCreation {
		Name: request.B.Name,
		Image: imageB,
		},
 	C: EntityCreation {
		Name: request.C.Name,
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
	TripleId TripleId  `parseArg:"triple_id,required"`
	Vote Vote `parseArg:"vote"`
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



type googleImageSearchResult struct {
	Items []struct{
		Image struct{
			ContextLink string
			ThumbnailLink string
		}
	}
}

type mfkImageSearchResult struct {
	Images []mfkImageSearchResultImage `json:"images"`
}

type mfkImageSearchResultImage struct {
	Context string `json:"context"`
	Url string `json:"url"`
}

// type googleImageSearchResultItem struct {
// 	Image googleImageSearchResultImage
// }

// type googleImageSearchResultImage struct {
// }

// type mfkImageSearch


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
	values.Set(		"filter", "1")
	values.Set(		"safe", "medium")
	values.Set(		"searchType", "image")
	values.Set(		"fields", "items(image(contextLink,thumbnailLink))")
	values.Set(		"pp", "1")
	values.Set(		"key", "AIzaSyDbjy0CKTMV5DoJR07ZYF5w-KL7Ey5lyGY")
	values.Set(		"q", query)
	values.Set(		"userIp", c.r.RemoteAddr)

	url := API_BASE + "?" + values.Encode()


	c.cxt.Infof("API Query: %s", url)

	fetcher := urlfetch.Client(c.cxt)
	response, err := fetcher.Get(url)
	if err != nil {
		panic(err)
	}

	var resultBuffer bytes.Buffer
	if _, err := resultBuffer.ReadFrom(response.Body); err != nil {
		panic(err);
	}

	var res googleImageSearchResult
	if err := json.Unmarshal(resultBuffer.Bytes(), &res); err != nil {
		panic(err)
	}

	var ourResult mfkImageSearchResult
	for _, item := range res.Items {
		img := item.Image
		ourResult.Images = append(ourResult.Images, mfkImageSearchResultImage{
		Context: img.ContextLink,
		Url: img.ThumbnailLink,
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
		panic(err);
	}
}
