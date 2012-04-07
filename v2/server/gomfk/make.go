package gomfk

import (
	"appengine"
	"appengine/datastore"
	_ "appengine/user"
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/url"
)

func makeHandler(w http.ResponseWriter, r *http.Request) {
	t := Templates()

	t.Execute(w, nil)
}

type makeRequestEntity struct {
	Name string
}

type makeRequest struct {
	A, B, C makeRequestEntity
}

func makeDoHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)

	var req makeRequest
	err := json.Unmarshal([]byte(r.FormValue("data")), &req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	triple := Triple{}
	triple.Init(NewRandom())

	triple.NameA = req.A.Name
	triple.NameB = req.B.Name
	triple.NameC = req.C.Name

	triple.Creator = UserIdFromContext(r)

	key := datastore.NewIncompleteKey(cxt, "Triple", nil)
	if key, err = datastore.Put(cxt, key, &triple); err != nil {
		cxt.Infof("%v", err)
	}

	fmt.Println(key.IntID())

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}
