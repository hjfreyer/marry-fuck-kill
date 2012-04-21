package json_api

import (
	"appengine"
	_ "appengine/datastore"
	_ "appengine/user"
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/url"
	_ "strings"
)

const RETRY_COUNT = 3

type ApiError struct {
	errorCode int
	message   string
}

func Error(errorCode int, msg string) *ApiError {
	return &ApiError{errorCode, msg}
}

func Error500(err error) *ApiError {
	return &ApiError{500, err.Error()}
}

type ApiMethod interface {
	Call(w http.ResponseWriter, r *http.Request)
}

type JsonApiMethod interface {
	Call(httpRequest *http.Request, request, response interface{}) *ApiError
	NewRequest() interface{}
	NewResponse() interface{}
}

type jsonAdapter struct {
	m JsonApiMethod
}

func (j jsonAdapter) Call(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	data := r.FormValue("data")

	requestMsg := j.m.NewRequest()
	responseMsg := j.m.NewResponse()

	// Parse request
	if err := json.Unmarshal([]byte(data), &requestMsg); err != nil {
		cxt.Errorf("Request data didn't parse. %v\n%s", err, data)
		http.Error(w, err.Error(), 400)
		return
	}

	if err := j.m.Call(r, requestMsg, responseMsg); err != nil {
		cxt.Errorf("Error %d while calling JSON method: %s",
			err.errorCode, err.message)
		http.Error(w, err.message, err.errorCode)
		return
	}

	response, err := json.Marshal(responseMsg)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}

func JsonMethod(m JsonApiMethod) ApiMethod {
	return jsonAdapter{m}
}

func NewApiHandler(methods map[string]ApiMethod) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cxt := appengine.NewContext(r)

		methodName := r.FormValue("method")
		cxt.Infof("Api called with method %q", methodName)

		if methodName == "" {
			http.Error(w, fmt.Sprintf("Must specify the 'method' parameter"), 400)
			return
		}

		method, ok := methods[methodName]
		if !ok {
			http.Error(w, fmt.Sprintf("Method %q not found.", methodName), 404)
			return
		}
		method.Call(w, r)
	}
}
