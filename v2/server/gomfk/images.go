package gomfk

import (
	"appengine"
	_ "appengine"
	_ "appengine/datastore"
	_ "appengine/datastore"
	_ "appengine/urlfetch"
	_ "appengine/user"
	_ "appengine/user"
	_ "encoding/json"
	_ "encoding/json"
	_ "fmt"
	_ "fmt"
	_ "gomfk/json_api"
	_ "html/template"
	_ "io/ioutil"
	_ "net/http"
	_ "net/http"
	_ "net/url"
	_ "net/url"
	_ "strings"
)

type ImageMessage struct {
	Url         string `json:"url"`
	ContentType string `json:"contentType"`
	SourceUrl   string `json:"sourceUrl"`
	Salt        int64  `json:"salt"`
	Hash        int64  `json:"hash"`
}

func StoreImage(cxt *appengine.Context,
	image ImageMessage) (id int64, err error) {
	fetcher := urlfetch.Client(cxt)

	response, err := fetcher.Get("https://www.googleapis.com/customsearch/v1?q=fox&cx=017343173679326196998%3Aomutomvh_wi&safe=medium&searchType=image&fields=items(image(thumbnailHeight%2CthumbnailLink%2CthumbnailWidth))&pp=1&key=AIzaSyDbjy0CKTMV5DoJR07ZYF5w-KL7Ey5lyGY")

	if err != nil {
		cxt.Infof("%s", err)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			cxt.Infof("%s", err)
			return
		}
		cxt.Infof("%s\n", string(contents))
	}
	return
}
