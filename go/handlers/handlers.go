package handlers

import (
	"appengine"
	"github.com/hjfreyer/marry-fuck-kill/go/mfklib"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
	"github.com/hjfreyer/marry-fuck-kill/go/impl"
	// _ "appengine/datastore"
	// 	"appengine/memcache"
_	"appengine/urlfetch"
	"appengine/user"
	// 	"bytes"
	"encoding/json"
_	"errors"
_	"fmt"
	// 	_ "gomfk/json_api"
_	"time"
	// 	"gomfk/parse_args"
	// 	"gomfk/mfklib"
	"io/ioutil"
	"net/http"
_	"net/url"
	// 	_ "reflect"
	// 	_ "strconv"
	// 	_ "strings"
//	"code.google.com/p/goprotobuf/proto"
	"regexp"
"strconv"
"html/template"
// 	_ "text/template"
)

// const RETRY_COUNT = 3

// func WrapHandler(handler func(c *Context)) func(
// 	w http.ResponseWriter, r *http.Request) {

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		cxt := appengine.NewContext(r)
// 		user := UserIdFromContext(r)

// 		c := Context{user, cxt, w, r}
// 		handler(&c)
// 	}
// }

// type Context struct {
// 	userId UserId
// 	cxt    appengine.Context
// 	w      http.ResponseWriter
// 	r      *http.Request
// }

// func (c *Context) Error(code int, err error) {
// 	c.Errorf(code, err.Error())
// }

// func (c *Context) Errorf(code int, f string, args ...interface{}) {
// 	msg := fmt.Sprintf(f, args...)
// 	c.cxt.Errorf(msg)
// 	http.Error(c.w, msg, code)
// }

// type makeRequest struct {
// 	A makeRequest_Entity `parseArg:"a"`
// 	B makeRequest_Entity `parseArg:"b"`
// 	C makeRequest_Entity `parseArg:"c"`
// }

// type makeRequest_Entity struct {
// 	Name  string            `parseArg:"Name,required"`
// 	Image makeRequest_Image `parseArg:"Image"`
// }

// type makeRequest_Image struct {
// 	Url string `parseArg:"Url,required"`
// 	// ContentType string `parseArg:"ContentType"`
// 	// SourceUrl   string `parseArg:"sourceUrl"`
// 	// Salt        int64  `parseArg:"salt"`
// 	// Hash        int64  `parseArg:"hash"`
// }

// type makeResponse struct {
// 	Id int64 `json:"id"`
// }



func MakeMFKImpl(req *http.Request) *mfklib.MFKImpl {
	cxt := appengine.NewContext(req)

	var userId string
	if u := user.Current(cxt); u != nil {
		userId = u.Email + "::" + u.ID
	}
	userId = req.RemoteAddr

	backend := impl.BackendImpl{cxt, req}
	db := impl.NewDb(cxt)

	return &mfklib.MFKImpl{
		UserId:        userId,
		Logger:        backend,
		ImageSearcher: backend,
		ImageFetcher: backend,
		Database: db,
	}
}

func checkOk(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseJsonRequest(r *http.Request, dest interface{}) error {
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, dest)
}

func ApiMakeHandler(w http.ResponseWriter, r *http.Request) {
	request := mfklib.MakeTripleRequest{}
	_ = ParseJsonRequest(r, &request)

	w.Write([]byte(proto.MarshalTextString(&request)))

	//	mfk := MakeMFKImpl(r)


}

type apiImageSearchHandler struct{}

func (a apiImageSearchHandler) GetKey(r *http.Request) string {
	return "IMGSRCH:" + r.FormValue("query")
}

func (a apiImageSearchHandler) Handle(r *http.Request) ([]byte, *Error) {
	query := r.FormValue("query")

	mfk := MakeMFKImpl(r)
	result, err := mfk.ImageSearch(query)
	checkOk(err)

	response, err := json.Marshal(result)
	checkOk(err)

	return response, nil
}

var ImageSearchApiHandler = NewCachedHandler(apiImageSearchHandler{})



type makeTripleHandler struct{}

func (m makeTripleHandler) Handle(w http.ResponseWriter, r *http.Request) *Error {
	if r.Method != "POST" {
		return &Error{http.StatusMethodNotAllowed, "Error: Use POST", nil}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	request := mfklib.MakeTripleRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		return &Error{400, "Request body is not a valid JSON MakeTripleRequest", err}
	}

	mfk := MakeMFKImpl(r)

	resp, err := mfk.MakeTriple(&request)
	respJson, err := json.Marshal(&resp)
	checkOk(err)
	w.Write(respJson)

	return nil
}

var MakeTripleApiHandler = NewErrorHandler(makeTripleHandler{})

func parseTripleId(t string) (mfklib.TripleId, *Error) {
	tripleId, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		if nerr := err.(*strconv.NumError); nerr.Err == strconv.ErrRange {
			return 0, &Error{404, "Triple ID too long", err}
		} else {
			panic(err)
		}
	}

	return mfklib.TripleId(tripleId), nil
}

func notFoundError(err error) *Error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case mfklib.EntityNotFoundError:
		return &Error{404, "Not found.", e}
	}
	panic(err)
}

type getImageHandler struct{}

var badUrlFormat = &Error{404, "Bad URL format", nil}

var IMAGE_RE, _ = regexp.Compile("^/i/([0-9]+)/([012])$")

func (getImageHandler) Handle(w http.ResponseWriter, r *http.Request) *Error{
	match := IMAGE_RE.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return badUrlFormat
	}

	tripleId, herr := parseTripleId(match[1])
	if herr != nil {
		return herr
	}
	entity := match[2]

	mfk := MakeMFKImpl(r)

	image, err := mfk.GetImage(tripleId, entity)
	if herr := notFoundError(err); herr != nil {
		return herr
	}

	w.Header().Set("content-type", *image.ContentType)
	_, err = w.Write(image.Data)
	checkOk(err)

	return nil
}

var GetImageHandler = NewErrorHandler(getImageHandler{})

// var request makeRequest
// err := parse_args.ParseArgs(c.r, &request)
// if err != nil {
// 	c.Error(400, err)
// 	return
// }

// 	FetchOrDie := func(url string) *FetchedImage {
// 		image, err := FetchImage(c.cxt, url)
// 		maybePanic(err)
// 		return image
// 	}

// 	imageA := FetchOrDie(request.A.Image.Url)
// 	imageB := FetchOrDie(request.B.Image.Url)
// 	imageC := FetchOrDie(request.C.Image.Url)

// 	triple := TripleCreation{
// 		A: EntityCreation{
// 			Name:  request.A.Name,
// 			Image: imageA,
// 		},
// 		B: EntityCreation{
// 			Name:  request.B.Name,
// 			Image: imageB,
// 		},
// 		C: EntityCreation{
// 			Name:  request.C.Name,
// 			Image: imageC,
// 		},
// 		Creator: c.userId,
// 	}

// 	db := NewAppengineDataAccessor(c.cxt)
// 	tripleId, err := db.MakeTriple(triple)
// 	if err != nil {
// 		panic(err)
// 	}

// 	response := makeResponse{int64(tripleId)}
// 	responseMsg, err := json.Marshal(response)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.w.Write(responseMsg)
// }

// type voteRequest struct {
// 	TripleId TripleId `parseArg:"triple_id,required"`
// 	Vote     Vote     `parseArg:"vote"`
// }

// func (v *voteRequest) Validate() error {
// 	if !v.Vote.IsValid() {
// 		return errors.New("Invalid vote: " + string(v.Vote))
// 	}
// 	return nil
// }

// func ApiVoteHandler(c *Context) {
// 	var request voteRequest
// 	err := parse_args.ParseArgs(c.r, &request)
// 	if err != nil {
// 		c.Error(400, err)
// 		return
// 	}

// 	db := NewAppengineDataAccessor(c.cxt)
// 	err = db.UpdateVote(request.TripleId, c.userId, request.Vote)

// 	if err != nil {
// 		panic(err)
// 	}
// }

// type HttpError struct {
// 	StatusCode int
// 	Message string
// }

// type Fetcher interface {
// 	Fetch(query string) *http.Response
// }

// func ImageSearch(query string, logger Logger, fetcher Fetcher) (ImageSearchResponse, HttpError) {
// 	var imageSearchJson []byte
// 	for attempt := 0; attempt < NUM_RETRIES; attempt++ {
// 		resp := fetcher.Fetch(query)
// 		defer resp.Body.Close()

// 		if resp.StatusCode == 200 {
// 			var err error
// 			imageSearchJson, err = ioutil.ReadAll(resp.Body)
// 			CheckOk(err)
// 		}
// 	}
// }

// func ApiImageSearchHandler(c *Context) {
// 	query := c.r.FormValue("query")
// 	c.cxt.Infof(query)

// 	if item, err := memcache.Get(c.cxt, query); err == memcache.ErrCacheMiss {
// 		c.cxt.Infof("Memcache: miss")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		c.cxt.Infof("Memcache: hit")
// 		c.w.Write(item.Value)
// 		return
// 	}

type tripleView struct {
	Id mfklib.TripleId
	*mfklib.Triple
}

func (tv tripleView) UserVoted() bool {
	return false
}

func (tv tripleView) Entities() [3]entityView {
	return [3]entityView{
		{*tv.A},
		{*tv.B},
		{*tv.C},
	}
}

type entityView struct {
	mfklib.Triple_Entity
}

func (ev entityView) ChartUrl(m, f, k int) string {
	return ""
}

type singleTripleHandler struct{}

var VOTE_RE, _ = regexp.Compile("^/vote/([0-9]+)$")

func (singleTripleHandler) Handle(w http.ResponseWriter, r *http.Request) *Error {
	match := VOTE_RE.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return badUrlFormat
	}

	tripleId, herr := parseTripleId(match[1])
	if herr != nil {
		return herr
	}

	// cxt := appengine.NewContext(r)


	mfk := MakeMFKImpl(r)


	triple, err := mfk.GetTriple(tripleId)
	if herr := notFoundError(err); herr != nil {
		return herr
	}
	// mfk

	// tripleIds, err := db.GetTripleIds(10)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// triples, err := db.GetTriples(user, tripleIds)
	// if err != nil {
	// 	cxt.Errorf("%v", []error(err.(appengine.MultiError)))
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// tripleViews := make([]TripleView, len(triples))
	// for i, t := range triples {
	// 	tripleViews[i] = TripleView(t)
	// }

	//	mfk.Infof("%s", triple)
	templates := Templates()
	if err := templates.ExecuteTemplate(w, "tripleList", []tripleView{{Id: tripleId,
		Triple: triple}}); err != nil {
		panic(err)
	}

	return nil
}

var SingleTripleHandler = NewErrorHandler(singleTripleHandler{})



func ListHandler(w http.ResponseWriter, r *http.Request) {
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	var makerStructure struct {
		Entities [3]struct {
			ResultBoxes [10]struct{}
		}
	}

	t := Templates()

	if err := t.ExecuteTemplate(w, "make", makerStructure); err != nil {
		panic(err)
	}
}

func Templates() *template.Template {
	return template.Must(template.ParseFiles(
		"templates/header.html",
		"templates/make.html",
		"templates/triple.html",
		"templates/vote.html",
		"generated_templates/js_include.html",
	))
}
