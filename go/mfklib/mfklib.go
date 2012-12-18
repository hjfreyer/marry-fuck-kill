package mfklib

import (
	"errors"
	"fmt"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
)

type MFKImpl struct {
	UserId

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
		CreatorId: proto.String(string(mfk.UserId)),
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
		Func:     "GetImage",
		Argument: "entity",
		Value:    entity,
		Cause:    "Must be 0, 1, or 2",
	}
}

func getIntsForVote(s *TripleStats, v VoteStatus) []*uint64 {
	switch v {
	case VoteStatus_UNSET:
		return []*uint64{}
	case VoteStatus_SKIP:
		return []*uint64{&s.Skips}
	case VoteStatus_MFK:
		return []*uint64{&s.A.Marry, &s.B.Fuck, &s.C.Kill}
	case VoteStatus_MKF:
		return []*uint64{&s.A.Marry, &s.B.Kill, &s.C.Fuck}
	case VoteStatus_FMK:
		return []*uint64{&s.A.Fuck, &s.B.Marry, &s.C.Kill}
	case VoteStatus_FKM:
		return []*uint64{&s.A.Fuck, &s.B.Kill, &s.C.Marry}
	case VoteStatus_KMF:
		return []*uint64{&s.A.Kill, &s.B.Marry, &s.C.Fuck}
	case VoteStatus_KFM:
		return []*uint64{&s.A.Kill, &s.B.Fuck, &s.C.Marry}
	}
	panic(fmt.Errorf("Invalid vote: %d", v))
}

func addVoteToStats(stats *TripleStats, vote VoteStatus) {
	for _, num := range getIntsForVote(stats, vote) {
		*num++
	}
}

func subtractVoteFromStats(stats *TripleStats, vote VoteStatus) {
	for _, num := range getIntsForVote(stats, vote) {
		*num--
	}
}

func (mfk MFKImpl) GetTripleStatsForUser(tripleId TripleId) (TripleStats, VoteStatus, error) {
	var	stats TripleStats
	var	vote VoteStatus

	err := mfk.UpdateStats(tripleId, mfk.UserId, &stats, &vote, func() (bool, error) {
		return false, nil
	})
	if err != nil {
		return stats, vote, err
	}

	subtractVoteFromStats(&stats, vote)

	return stats, vote, nil
}

func (mfk MFKImpl) ChangeVote(tripleId TripleId, newVote VoteStatus) error {
	var stats TripleStats
	var prevVote VoteStatus

	err := mfk.UpdateStats(tripleId, mfk.UserId, &stats, &prevVote, func() (bool, error) {
		if newVote == prevVote {
			return false, nil
		}
		subtractVoteFromStats(&stats, prevVote)
		addVoteToStats(&stats, newVote)
		prevVote = newVote
		return true, nil
	})
	if err != nil {
		return err
	}
	return nil
}