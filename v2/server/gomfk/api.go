package gomfk

import (
	"appengine"
	_"appengine/datastore"
	_ "appengine/user"
	_"encoding/json"
	_ "fmt"
	"gomfk/json_api"
	"net/http"
	_ "net/url"
	_ "strings"
)

const RETRY_COUNT = 3

type apiError struct {
	errorCode int
	errorStr  string
}

func (a *apiError) Error() string {
	return a.errorStr
}

var ApiHandler = json_api.NewApiHandler(map[string]json_api.ApiMethod{
	"make": json_api.JsonMethod(MakeMethod{}),
	"vote": json_api.JsonMethod(VoteMethod{}),
})

type makeRequest struct {
	A makeRequest_Entity
	B makeRequest_Entity
	C makeRequest_Entity
}

type makeRequest_Entity struct {
	Name  string
	Image ImageMessage
}

type makeResponse struct {
	Id int64 `json:"id"`
}

type MakeMethod struct{}
func (m MakeMethod) NewRequest() interface{}  { return &makeRequest{} }
func (m MakeMethod) NewResponse() interface{} { return &makeResponse{} }
func (m MakeMethod) Call(httpRequest *http.Request,
	iRequest, iResponse interface{}) *json_api.ApiError {
	cxt := appengine.NewContext(httpRequest)
	request, response := iRequest.(*makeRequest), iResponse.(*makeResponse)

	if request.A.Name == "" || request.B.Name == "" || request.C.Name == "" {
		return json_api.Error(400, "Field missing")
	}

	imageA, err := FetchImage(cxt, request.A.Image)
	if err != nil {
		return json_api.Error500(err)
	}
	imageB, err := FetchImage(cxt, request.B.Image)
	if err != nil {
		return json_api.Error500(err)
	}
	imageC, err := FetchImage(cxt, request.C.Image)
	if err != nil {
		return json_api.Error500(err)
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
	Creator: UserIdFromContext(httpRequest),
	}

 	db := NewAppengineDataAccessor(cxt)
	tripleId, err := db.MakeTriple(triple)
	if err != nil {
		return json_api.Error500(err)
	}

	response.Id = int64(tripleId)
	return nil
}

type voteRequest struct {
	TripleId              int64  `json:"triple_id"`
	Vote string `json:"vote"`
}

type voteResponse struct {}

type VoteMethod struct{}
func (m VoteMethod) NewRequest() interface{}  { return &voteRequest{} }
func (m VoteMethod) NewResponse() interface{} { return &voteResponse{} }
func (m VoteMethod) Call(httpRequest *http.Request,
	iRequest, iResponse interface{}) *json_api.ApiError {
	request := iRequest.(*voteRequest)

	vote := Vote(request.Vote)

	if !vote.IsValid() {
		return json_api.Error(400, "Invalid vote " + string(vote))
	}

	cxt := appengine.NewContext(httpRequest)
	user := UserIdFromContext(httpRequest)

 	db := NewAppengineDataAccessor(cxt)
	err := db.UpdateVote(TripleId(request.TripleId), user, vote)

	if err != nil {
		return json_api.Error500(err)
	}

	return nil
}
