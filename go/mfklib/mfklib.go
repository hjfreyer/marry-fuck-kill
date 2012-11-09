
package mfklib

import (
	"net/http"
)

type Error struct {
	Error error
	Message string
	StatusCode int
}

type Logger interface {
    Infof(format string, args ...interface{})
    Warningf(format string, args ...interface{})
}

type ImageSearchResult struct {
	ImageUrl string
	ContextUrl string
}

type ImageSearcher func(query string) (*http.Response, err)

func ImageSearch(query string, log Logger, search ImageSearcher, retries int) (
	ImageSearchResponse, *Error) {

	json, err := func() ([]byte, *Error) {
		for attempt := 0; attempt < retries; attempt++ {
			resp, err := search(query)
			if err != nil {
				log.Warningf("%s", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				json, err = ioutil.ReadAll(resp.Body)
				CheckOk(err)
				return
			}
		}
	}()
}