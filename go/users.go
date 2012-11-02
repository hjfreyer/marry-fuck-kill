package gomfk

import (
	"appengine"
	"appengine/user"
	_ "fmt"
	"net/http"
)

type UserId string

func UserIdFromContext(req *http.Request) UserId {
	cxt := appengine.NewContext(req)
	u := user.Current(cxt)

	if u != nil {
		return UserId(u.Email + "::" + u.ID)
	}
	return UserId(req.RemoteAddr)
}
