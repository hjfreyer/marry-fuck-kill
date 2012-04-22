package gomfk

import (
	"appengine"
	_ "appengine"
	"appengine/datastore"
	_ "appengine/datastore"
	 "appengine/urlfetch"
	_ "appengine/user"
	_ "appengine/user"
	_ "encoding/json"
	_ "encoding/json"
	_ "fmt"
	_ "fmt"
	_ "gomfk/json_api"
	_ "html/template"
	 "io/ioutil"
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

func StoreImage(cxt appengine.Context, image ImageMessage) (int64, error) {
	// TODO(hjfreyer): Verify image.
	fetcher := urlfetch.Client(cxt)

	cxt.Infof("Attempting to fetch url: %s", image.Url)
	response, err := fetcher.Get(image.Url)
	if err != nil {
		return 0, err
	}

	if response.StatusCode != 200 {
		panic(response)
	}

	contentType := response.Header.Get("content-type")
	cxt.Infof("Content type: %s", contentType)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	// TODO(hjfreyer): Check that the data isn't too large
	entityImage := EntityImage{
	SourceUrl : image.SourceUrl,
	ContentType : contentType,
	Data : contents,
	}

	key := datastore.NewIncompleteKey(cxt, "EntityImage", nil)
	key, err = datastore.Put(cxt, key, &entityImage)
	if err != nil {
		return 0, err
	}

	// cxt.Infof("%s\n", string(contents))
	return key.IntID(), nil
}
