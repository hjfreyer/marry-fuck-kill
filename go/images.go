package gomfk

import (
	"appengine"
	_ "appengine"
	_ "appengine/datastore"
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

type FetchedImage struct {
	SourceUrl   string
	ContentType string
	Data        []byte
}

func FetchImage(cxt appengine.Context, url string) (
	*FetchedImage, error) {
	// TODO(hjfreyer): Verify image.
	fetcher := urlfetch.Client(cxt)

	cxt.Infof("Attempting to fetch url: %s", url)
	response, err := fetcher.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		panic(response)
	}

	contentType := response.Header.Get("content-type")
	cxt.Infof("Content type: %s", contentType)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// TODO(hjfreyer): Check that the data isn't too large

	return &FetchedImage{
		ContentType: contentType,
		Data:        contents,
	}, nil
}
