package mfklib

import (
	"code.google.com/p/goprotobuf/proto"
)

type Error struct {
	StatusCode int
	Message    string
	Error      error
}

type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
}

type ImageSearcher interface {
	Search(query string) ([]ImageMetadata, error)
}

type ImageOrError struct {
	Image
	error
}

type ImageFetcher interface {
	FetchImage(url string) chan ImageOrError
}

type TripleStatsUpdater func(*TripleStats)

type Database interface {
	AddTriple(*Triple) (int64, error)
	GetTriple(triple_id int64) (*Triple, error)

	UpdateStats(triple_id int64, updater TripleStatsUpdater) error
}

type MFKImpl struct {
	UserId string

	Logger
	ImageSearcher
	ImageFetcher
	Database
}

func (mfk MFKImpl) ImageSearch(query string) (*ImageSearchResponse, *Error) {
	results, err := mfk.Search(query)

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

func (mfk MFKImpl) MakeTriple(request *MakeTripleRequest) (*MakeTripleResponse, *Error) {
	if !mfk.VerifyWrappedImage(request.A.Image) ||
		!mfk.VerifyWrappedImage(request.B.Image) ||
		!mfk.VerifyWrappedImage(request.C.Image) {
		return nil, &Error{400, "Image verification failed.", nil}
	}

	fetchA := mfk.FetchImage(*request.A.Image.Metadata.Url)
	fetchB := mfk.FetchImage(*request.B.Image.Metadata.Url)
	fetchC := mfk.FetchImage(*request.C.Image.Metadata.Url)

	imgA := <-fetchA
	imgB := <-fetchB
	imgC := <-fetchC

	for _, img := range []ImageOrError{imgA, imgB, imgC} {
		if img.error != nil {
			return nil, &Error{500, "Error fetching image.", img.error}
		}
	}

	triple := Triple{
		CreatorId: proto.String(mfk.UserId),
		A: &Triple_Entity{
			Name:  request.A.Name,
			Image: &imgA.Image,
		},
	}

	tripleId, err := mfk.AddTriple(&triple)

	if err != nil {
		return nil, &Error{500, "Error storing triple.", err}
	}

	response := MakeTripleResponse{
		TripleId: proto.Int64(tripleId),
	}

	return &response, nil
}

func (mfk MFKImpl) VerifyWrappedImage(image *WrappedImageMetadata) bool {
	return true
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
