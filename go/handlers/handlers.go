package handlers

import (
	"appengine"
	"github.com/hjfreyer/marry-fuck-kill/go/impl"
	"github.com/hjfreyer/marry-fuck-kill/go/mfklib"
	"appengine/user"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"strconv"
)

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
		ImageFetcher:  backend,
		Database:      db,
	}
}

func checkOk(err error) {
	if err != nil {
		panic(err)
	}
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

func (getImageHandler) Handle(w http.ResponseWriter, r *http.Request) *Error {
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

type tripleView struct {
	*mfklib.Triple

	Id mfklib.TripleId

	MarryCount uint64
	FuckCount  uint64
	KillCount  uint64
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
