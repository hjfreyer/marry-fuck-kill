package mfklib

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
)

type EntityNotFoundError struct {
	Type string
	Id int64
	Err error
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("Entity of type %q with id %d not found. Error: %s",
		e.Type, e.Id, e.Err)
}

type IllegalArgumentError struct {
	Func string
	Argument string
	Value interface{}
	Cause string
}

func (e IllegalArgumentError) Error() string {
	return fmt.Sprintf("Function %q got illegal value %q for argument %q: %s",
		e.Func, e.Value, e.Argument, e.Cause)
}

type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
}

type ImageSearcher interface {
	Search(query string) ([]*ImageMetadata, error)
}

type ImageOrError struct {
	*Image
	error
}

type ImageFetcher interface {
	FetchImage(metadata *ImageMetadata) chan ImageOrError
}

type TripleStatsUpdater func(*TripleStats, Vote_VoteType) Vote_VoteType

type Database interface {
	AddTriple(*Triple) (TripleId, error)
	GetTriple(tripleId TripleId) (*Triple, error)

	UpdateStats(tripleId TripleId, updater TripleStatsUpdater) error
}

type TripleId int64

type MFKImpl struct {
	UserId string

	Logger
	ImageSearcher
	ImageFetcher
	Database
}

func (mfk MFKImpl) ImageSearch(query string) (*ImageSearchResponse, error) {
	results, err := mfk.Search(query)

	if err != nil {
		return nil, err
	}

	response := ImageSearchResponse{}
	for _, img := range results {
		response.Image = append(response.Image, &WrappedImageMetadata{
			Metadata: img,
		})
	}

	return &response, nil
}

var ImageVerificationFailedError = errors.New("Image verification failed.")

func (mfk MFKImpl) MakeTriple(request *MakeTripleRequest) (*MakeTripleResponse, error) {
	if !mfk.verifyWrappedImage(request.A.Image) ||
		!mfk.verifyWrappedImage(request.B.Image) ||
		!mfk.verifyWrappedImage(request.C.Image) {
		return nil, ImageVerificationFailedError
	}

	fetchA := mfk.FetchImage(request.A.Image.Metadata)
	fetchB := mfk.FetchImage(request.B.Image.Metadata)
	fetchC := mfk.FetchImage(request.C.Image.Metadata)

	imgA := <-fetchA
	imgB := <-fetchB
	imgC := <-fetchC

	for _, img := range []ImageOrError{imgA, imgB, imgC} {
		if img.error != nil {
			return nil, img.error
		}
	}

	triple := Triple{
		CreatorId: proto.String(mfk.UserId),
		A: &Triple_Entity{
			Name:  request.A.Name,
			Image: imgA.Image,
		},
		B: &Triple_Entity{
			Name:  request.B.Name,
			Image: imgB.Image,
		},
		C: &Triple_Entity{
			Name:  request.C.Name,
			Image: imgC.Image,
		},
	}

	tripleId, err := mfk.AddTriple(&triple)

	if err != nil {
		return nil, err
	}

	response := MakeTripleResponse{
		TripleId: proto.Int64(int64(tripleId)),
	}

	return &response, nil
}

func (mfk MFKImpl) verifyWrappedImage(image *WrappedImageMetadata) bool {
	return true
}

func (mfk MFKImpl) GetImage(tripleId TripleId, entity string) (*Image, error) {
	triple, err := mfk.GetTriple(tripleId)
	if err != nil {
		return nil, err
	}
	switch entity {
	case "0":
		return triple.A.Image, nil
	case "1":
		return triple.B.Image, nil
	case "2":
		return triple.C.Image, nil
	}
	return nil, &IllegalArgumentError{
		Func: "GetImage",
		Argument: "entity",
		Value: entity,
		Cause: "Must be 0, 1, or 2",
	}
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
