package gomfk

import (
_	"strings"
	"appengine"
	"appengine/datastore"
	_ "appengine/user"
	"encoding/json"
_	"fmt"
	"net/http"
	_ "net/url"
	"gomfk/json_api"
)

const RETRY_COUNT = 3

type apiError struct {
	errorCode int
	errorStr string
}

func (a *apiError) Error() string {
	return a.errorStr
}

var ApiHandler = json_api.NewApiHandler(map[string]json_api.ApiMethod{
	"make" : json_api.JsonMethod(MakeMethod{}),
})

type makeRequest struct {
	A makeRequest_Entity
	B makeRequest_Entity
	C makeRequest_Entity
}

type makeRequest_Entity struct {
	Name string
	Image ImageMessage
}

type makeResponse struct {
	Id int64 `json:"id"`
}

type MakeMethod struct{}

func (m MakeMethod) NewRequest() interface{} { return &makeRequest{} }
func (m MakeMethod) NewResponse() interface{} { return &makeResponse{} }
func (m MakeMethod) Call(httpRequest *http.Request,
	iRequest, iResponse interface{}) *json_api.ApiError {
	cxt := appengine.NewContext(httpRequest)
	request, response := iRequest.(*makeRequest), iResponse.(*makeResponse)

	if request.A.Name == "" || request.B.Name == "" || request.C.Name == "" {
		return json_api.Error(400, "Field missing")
	}

	triple := Triple{}
	triple.Init(NewRandom())

	triple.NameA = request.A.Name
	var err error
	triple.ImageIdA, err = StoreImage(request.A.Image)
	if err != nil { return json_api.Error500(err) }

	triple.NameB = request.B.Name
	triple.ImageIdA, err = StoreImage(request.A.Image)
	if err != nil { return json_api.Error500(err) }

	triple.NameC = request.C.Name
	triple.ImageIdA, err = StoreImage(request.A.Image)
	if err != nil { return json_api.Error500(err) }

	triple.Creator = UserIdFromContext(httpRequest)

	key := datastore.NewIncompleteKey(cxt, "Triple", nil)
	key, err = datastore.Put(cxt, key, &triple)
	if err != nil { return json_api.Error500(err) }

	response.Id = key.IntID()
	return nil
}


type voteRequest struct {
	Triple_ID int64
	Vote_A, Vote_B, Vote_C string
}

func VoteApiHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)

	var req voteRequest
	req.Triple_ID = -1
	if err := json.Unmarshal([]byte(r.FormValue("data")), &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if req.Triple_ID == -1 {
		http.Error(w, "Must specify triple_id", 400)
		return
	}

	voteSet := make(map[byte]bool)
	voteSet[req.Vote_A[0]] = true
	voteSet[req.Vote_B[0]] = true
	voteSet[req.Vote_C[0]] = true

	if !voteSet['m'] || !voteSet['f'] || !voteSet['k'] {
		http.Error(w, "Votes must be m, f, and k", 400)
		return
	}

	newVote := Vote{[]int{
		VoteFromChar(req.Vote_A[0]),
		VoteFromChar(req.Vote_B[0]),
		VoteFromChar(req.Vote_C[0]),
	}}

	user := UserIdFromContext(r)
	cxt.Infof("User: %v", user)

	tripleKey := datastore.NewKey(cxt, "Triple", "", req.Triple_ID, nil)
	voteKey := datastore.NewKey(cxt, "Vote", string(user), 0, tripleKey)


	for tryTime := 0; tryTime < RETRY_COUNT; tryTime++ {
		triple := new(Triple)
		oldVote := new(Vote)
		err := datastore.RunInTransaction(cxt, func(c appengine.Context) error {
			if err := datastore.Get(cxt, tripleKey, triple); err != nil {
				return err
			}

			err := datastore.Get(cxt, voteKey, oldVote)
			if err != nil && err != datastore.ErrNoSuchEntity {
				return err
			}
			if err != datastore.ErrNoSuchEntity {
				triple.SubtractVote(*oldVote)
			}

			triple.AddVote(newVote)

			if _, err := datastore.PutMulti(cxt, []*datastore.Key{tripleKey, voteKey},
				[]interface{}{triple, &newVote}); err != nil {
				return err
			}

			return nil
		}, nil)

		if err == nil {
			cxt.Infof("Transaction succeeded")
			return
		} else if err == datastore.ErrNoSuchEntity {
			http.Error(w, "No such entity.", 404)
			return
		} else {
			cxt.Errorf("Transaction failed: %v", err)
		}
	}
	http.Error(w, "Failed to record vote.", 500)
}
