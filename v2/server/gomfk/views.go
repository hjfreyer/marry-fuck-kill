package gomfk

import (

"strings"
"strconv"
	"appengine"
	"appengine/datastore"
	_ "appengine/urlfetch"
	_ "appengine/user"
	_ "encoding/json"
	_ "fmt"
	"html/template"
	_ "io/ioutil"
	"net/http"
	_ "net/url"
)

func Templates() *template.Template {
	return template.Must(template.ParseFiles(
		"gomfk/templates/header.html",
		"gomfk/templates/make.html",
		"gomfk/templates/triple.html",
		"gomfk/templates/vote.html",
	))
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, "404 not found!!!1", 404)
}

// Transport for request that need an API key
// This is meant to be used like this:
//     t := &googleapi.Transport{ApiKey: config.key,
//         http.RoundTripper: http.DefaultTransport}
//     client := &http.Client{Transport: t}
type Transport struct {
	ApiKey string
	http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	args := req.URL.Query()
	args.Set("key", t.ApiKey)
	req.URL.RawQuery = args.Encode()
	r, _ := http.NewRequest(req.Method, req.URL.String(), req.Body)
	return t.RoundTripper.RoundTrip(r)
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)

	imageId, err := strconv.ParseInt(strings.Replace(r.URL.Path, "/i/", "", 1), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cxt.Infof("Looking up image: %d", imageId)

	key := datastore.NewKey(cxt, "EntityImage", "", imageId, nil)

	var image EntityImage
	if err := datastore.Get(cxt, key, &image); err != nil {
		// TODO(hjfreyer): standardize
		http.Error(w, "Image not found", 404)
		return
	}

	w.Header().Set("content-type", image.ContentType)
	w.Write(image.Data)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFound(w)
		return
	}
	cxt := appengine.NewContext(r)

	// fetcher := urlfetch.Client(cxt)

	// response, err := fetcher.Get("https://www.googleapis.com/customsearch/v1?q=fox&cx=017343173679326196998%3Aomutomvh_wi&safe=medium&searchType=image&fields=items(image(thumbnailHeight%2CthumbnailLink%2CthumbnailWidth))&pp=1&key=AIzaSyDbjy0CKTMV5DoJR07ZYF5w-KL7Ey5lyGY")

	// if err != nil {
	//   cxt.Infof("%s", err)
	// 	return
	// } else {
	//   defer response.Body.Close()
	//   contents, err := ioutil.ReadAll(response.Body)
	//   if err != nil {
	//     cxt.Infof("%s", err)
	// 		return
	//   }
	//   cxt.Infof("%s\n", string(contents))
	// }

	user := UserIdFromContext(r)
	cxt.Infof("User: %v", user)

	q := datastore.NewQuery("Triple").Limit(100)
	templates := Templates()

	triples := make([]Triple, 0, 10)
	if _, err := q.GetAll(cxt, &triples); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cxt.Infof("%v", triples)

	if err := templates.ExecuteTemplate(w, "tripleList", triples); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	t := Templates()

	t.ExecuteTemplate(w, "make", nil)
}

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	t := Templates()

	t.ExecuteTemplate(w, "vote", nil)
}
