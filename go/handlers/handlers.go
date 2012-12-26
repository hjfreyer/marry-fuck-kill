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
	"strings"
)

func MakeMFKImpl(req *http.Request) *mfklib.MFKImpl {
	cxt := appengine.NewContext(req)

	var userId string
	if u := user.Current(cxt); u != nil {
		userId = u.Email + "::" + u.ID
	} else {
		userId = req.RemoteAddr
	}
	backend := impl.BackendImpl{cxt, req}
	db := impl.NewDb(cxt)

	return &mfklib.MFKImpl{
		UserId:        mfklib.UserId(userId),
		Logger:        backend,
		ImageSearcher: backend,
		ImageFetcher:  backend,
		Database:      db,
	}
}

func panicOnError(err error) {
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
	panicOnError(err)

	response, err := json.Marshal(result)
	panicOnError(err)

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
	panicOnError(err)
	w.Write(respJson)

	return nil
}

var MakeTripleApiHandler = NewErrorHandler(makeTripleHandler{})

type voteApiHandler struct{}

func (voteApiHandler) Handle(w http.ResponseWriter, r *http.Request) *Error {
	// if r.Method != "POST" {
	// 	return NewError(http.StatusMethodNotAllowed, nil, "GET not allowed - use POST")
	// }

	tripleIdStr := r.FormValue("triple_id")
	tripleId, herr := parseTripleId(tripleIdStr)
	if herr != nil {
		return herr
	}

	voteStr := strings.ToUpper(r.FormValue("vote"))
	vote, ok := mfklib.VoteStatus_value[voteStr]
	if !ok {
		return NewError(400, nil, "Invalid vote: %q", voteStr)
	}

	mfk := MakeMFKImpl(r)

	err := mfk.ChangeVote(tripleId, mfklib.VoteStatus(vote))
	if err != nil {
		return tripleNotFound(err)
	}
	return nil
}

var VoteApiHandler = NewErrorHandler(voteApiHandler{})

func parseTripleId(t string) (mfklib.TripleId, *Error) {
	tripleId, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		switch err.(*strconv.NumError).Err {
		case strconv.ErrSyntax:
			return mfklib.TripleId(0), NewError(400, err, "Not a valid Triple ID: %q", t)
		case strconv.ErrRange:
			return mfklib.TripleId(0), NewError(400, err, "Triple ID too large: %q", t)
		}
	}

	return mfklib.TripleId(tripleId), nil
}

func tripleNotFound(err error) *Error {
	tnf := err.(*mfklib.TripleNotFoundError)
	return NewError(404, err, "No Triple with ID %d", tnf.TripleId)
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

	// Regex ensures this is correct.
	entity := match[2]

	mfk := MakeMFKImpl(r)

	image, err := mfk.GetImage(tripleId, entity)
	if err != nil {
		return tripleNotFound(err)
	}

	w.Header().Set("content-type", *image.ContentType)
	_, err = w.Write(image.Data)
	panicOnError(err)

	return nil
}

var GetImageHandler = NewErrorHandler(getImageHandler{})

type tripleView struct {
	*mfklib.Triple

	Id mfklib.TripleId

	Stats mfklib.TripleStats
	Vote mfklib.VoteStatus
}

func (tv tripleView) UserVoted() bool {
	return !(tv.Vote == mfklib.VoteStatus_UNSET || tv.Vote == mfklib.VoteStatus_SKIP)
}

func (tv tripleView) Entities() [3]entityView {
	return [3]entityView{
		{tv.A, tv.Stats.A},
		{tv.B, tv.Stats.B},
		{tv.C, tv.Stats.C},
	}
}

func (tv tripleView) VotesPerEntity() [3]string {
	switch tv.Vote {
	case mfklib.VoteStatus_MFK:
		return [3]string{"marry", "fuck", "kill"}
	case mfklib.VoteStatus_MKF:
		return [3]string{"marry", "kill", "fuck"}
	case mfklib.VoteStatus_FMK:
		return [3]string{"fuck", "marry", "kill"}
	case mfklib.VoteStatus_FKM:
		return [3]string{"fuck", "kill", "marry"}
	case mfklib.VoteStatus_KMF:
		return [3]string{"kill", "marry", "fuck"}
	case mfklib.VoteStatus_KFM:
		return [3]string{"kill", "fuck", "marry"}
	}
	panic("Invalid vote")
}

type entityView struct {
	*mfklib.Triple_Entity
	mfklib.Tally
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

	cxt := appengine.NewContext(r)

	mfk := MakeMFKImpl(r)

	triple, err := mfk.GetTriple(tripleId)
	if err != nil {
		return tripleNotFound(err)
	}

	stats, vote, err := mfk.GetTripleStatsForUser(tripleId)
	panicOnError(err)

	tv := tripleView{
		Id: tripleId,
		Triple: triple,
		Stats: stats,
		Vote: vote,
	}

	login, _ := user.LoginURL(cxt, r.URL.Path)
	context := map[string]interface{}{
		"Triple" : tv,
		"LoginUrl" : login,
	}
	templates := Templates()
	panicOnError(templates.ExecuteTemplate(w, "singleTriplePage", context))

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
