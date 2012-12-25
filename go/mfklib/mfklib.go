package mfklib

import (
	"errors"
	"fmt"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

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

func (mfk MFKImpl) GetTriple(tripleId TripleId) (*Triple, error) {
	triple, err := mfk.Database.GetTriple(tripleId)
	switch err.(type) {
	case nil:
		return triple, nil
	case *TripleNotFoundError:
		return nil, err
	}
	panic(err)
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
	panic(fmt.Errorf("Invalid entity %q when getting image for Triple %d", entity, tripleId))
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

func (mfk MFKImpl) GetTripleStatsForUser(tripleId TripleId) (
	stats TripleStats, vote VoteStatus, err error) {
	if _, err = mfk.GetTriple(tripleId); err != nil {
		return // Returns err
	}

	err = mfk.UpdateStats(tripleId, mfk.UserId, &stats, &vote, func() bool {
		return false
	})
	panicOnError(err)

	subtractVoteFromStats(&stats, vote)

	return
}

func (mfk MFKImpl) ChangeVote(tripleId TripleId, newVote VoteStatus) error {
	if _, err := mfk.GetTriple(tripleId); err != nil {
		return err
	}

	var stats TripleStats
	var prevVote VoteStatus

	err := mfk.UpdateStats(tripleId, mfk.UserId, &stats, &prevVote, func() bool {
		if newVote == prevVote {
			return false
		}
		subtractVoteFromStats(&stats, prevVote)
		addVoteToStats(&stats, newVote)
		prevVote = newVote
		return true
	})
	panicOnError(err)
	return nil
}