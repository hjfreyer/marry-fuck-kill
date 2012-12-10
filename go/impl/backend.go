package impl

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"github.com/hjfreyer/marry-fuck-kill/go/mfklib"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
	"io/ioutil"
	"net/http"
	"net/url"
)

type BackendImpl struct {
	appengine.Context
	*http.Request
}

func (b BackendImpl) Search(query string) ([]*mfklib.ImageMetadata, error) {
	b.Infof("Searching for query '%s'", query)

	const API_BASE = "https://www.googleapis.com/customsearch/v1"

	values := url.Values{}
	values.Set("cx", "017343173679326196998:omutomvh_wi")
	values.Set("filter", "1")
	values.Set("safe", "medium")
	values.Set("searchType", "image")
	values.Set("fields", "items(image(contextLink,thumbnailLink))")
	values.Set("pp", "1")
	values.Set("key", "AIzaSyDbjy0CKTMV5DoJR07ZYF5w-KL7Ey5lyGY")
	values.Set("q", query)
	values.Set("userIp", b.RemoteAddr)

	url := API_BASE + "?" + values.Encode()
	b.Infof("API Query: %s", url)

	response, err := urlfetch.Client(b.Context).Get(url)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response: %d\n%s", response.StatusCode, body)
	}

	var jsonResult struct {
		Items []struct {
			Image struct {
				ContextLink   string
				ThumbnailLink string
			}
		}
	}
	if err := json.Unmarshal(body, &jsonResult); err != nil {
		panic(err)
	}

	result := make([]*mfklib.ImageMetadata, 0)
	for _, item := range jsonResult.Items {
		result = append(result, &mfklib.ImageMetadata{
			Context: proto.String(item.Image.ContextLink),
			Url:     proto.String(item.Image.ThumbnailLink),
		})
	}
	return result, nil
}

func (b BackendImpl) FetchImage(metadata *mfklib.ImageMetadata) chan mfklib.ImageOrError {
	fetch := func() (*mfklib.Image, error) {
		response, err := urlfetch.Client(b.Context).Get(*metadata.Url)
		defer response.Body.Close()
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		if response.StatusCode != 200 {
			return nil, fmt.Errorf("Bad response: %d\n%s", response.StatusCode, body)
		}

		return &mfklib.Image{
			Metadata:    metadata,
			ContentType: proto.String(response.Header.Get("Content-Type")),
			Data:        body,
		}, nil
	}

	out := make(chan mfklib.ImageOrError, 1)
	go func() {
		image, err := fetch()
		out <- mfklib.ImageOrError{image, err}
	}()
	return out
}
