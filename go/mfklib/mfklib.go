
package mfklib

import (
_	"net/http"
_	"io/ioutil"
_	"encoding/json"
)

type Error struct {
	StatusCode int
	Message string
	Error error
}

type Logger interface {
    Infof(format string, args ...interface{})
    Warningf(format string, args ...interface{})
}

type ImageSearcher interface {
	Search(query string) ([]ImageMetadata, error)
}

type MFKImpl struct {
	Log Logger
	ImageSearcher ImageSearcher
}

func (mfk MFKImpl) ImageSearch(query string) (*ImageSearchResponse, *Error) {
	results, err := mfk.ImageSearcher.Search(query)

	if err != nil {
		return nil, &Error{500, "Could not get image search results.", nil}
	}

	response := ImageSearchResponse{}
	for _, img := range results {
		response.Image = append(response.Image, &WrappedImageMetadata{
			Metadata: &img,
		})
	}

	return &response, nil
}




	// var imageSearchJson struct {
	// 	Items []struct {
	// 		Image struct {
	// 			ContextLink   string
	// 			ThumbnailLink string
	// 		}
	// 	}
	// }

	// tryImageSearch := func() bool {
	// 	resp, err := mfk.ImageSearcher.Search(query)
	// 	defer resp.Body.Close()

	// 	if err != nil {
	// 		mfk.Log.Warningf("%s", err)
	// 		return false
	// 	}

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	CheckOk(err)
	// 	if resp.StatusCode != 200 {
	// 		mfk.Log.Warningf("Image search returned error: %s\n%s", resp.Status, body)
	// 		return false
	// 	}

	// 	if json.Unmarshal(body, &imageSearchJson); err != nil {
	// 		mfk.Log.Warningf("Failed to parse JSON:\n%s", body)
	// 		return false
	// 	}
	// 	return true
	// }

	// success := false
	// for attempt := 0; attempt < retries; attempt++ {
	// 	mfk.Log.Infof("Attempt number %d to search for '%s'", attempt, query)
	// 	if tryImageSearch() {
	// 		mfk.Log.Infof("Attempt number %d to search for '%s'", attempt, query)
	// 		success = true
	// 		break;
	// 	}
	// }
	// if !success {
	// 	return nil, &Error{500, "Could not get image search results.", nil}
	// }

	// response := ImageSearchResponse{}
	// for _, item := range imageSearchJson.Items {
	// 	img := item.Image
	// 	response.Image = append(response.Image, &ImageMetadata{
	// 		ContextUrl: &img.ContextLink,
	// 		Url:     &img.ThumbnailLink,
	// 	})
	// }

	// return &response, nil
