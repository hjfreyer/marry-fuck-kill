package mfklib

import (
	"errors"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
	. "launchpad.net/gocheck"
	"testing"
	_"fmt"
)

type protoEqualsChecker struct {
	*CheckerInfo
}

// The ProtoEquals checker verifies that the obtained and expected values are
// both Google protocol buffers, and that they are proto-equal.
var ProtoEquals Checker = &protoEqualsChecker{
	&CheckerInfo{Name: "ProtoEquals", Params: []string{"obtained", "expected"}},
}

func (checker *protoEqualsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained, ok := params[0].(proto.Message)
	if !ok {
		return false, "obtained is not a proto.Message"
	}

	expected, ok := params[1].(proto.Message)
	if !ok {
		return false, "expected is not a proto.Message"
	}

	params[0] = proto.MarshalTextString(obtained)
	params[1] = proto.MarshalTextString(expected)

	return proto.Equal(obtained, expected), ""
}

func Test(t *testing.T) { TestingT(t) }

type TestLogger struct{ *C }

func (t TestLogger) Infof(format string, args ...interface{}) {
	t.Logf(format, args...)
}

func (t TestLogger) Warningf(format string, args ...interface{}) {
	t.Logf("W: "+format, args...)
}

type FakeImageSearcher func(query string) (results []*ImageMetadata, err error)
type FakeImageFetcher func(metadata *ImageMetadata) chan ImageOrError

func (f FakeImageSearcher) Search(query string) (results []*ImageMetadata, err error) {
	return f(query)
}

func (f FakeImageFetcher) FetchImage(metadata *ImageMetadata) chan ImageOrError {
	return f(metadata)
}

type tripleUserIdPair struct {
	TripleId
	UserId
}

type inMemoryDatabase struct {
	triples []Triple
	stats   []TripleStats
	vote  map[tripleUserIdPair]VoteStatus
}

func NewInMemoryDb() *inMemoryDatabase {
	return &inMemoryDatabase{
		vote: make(map[tripleUserIdPair]VoteStatus),
	}
}

func (db *inMemoryDatabase) AddTriple(triple *Triple) (TripleId, error) {
	id := len(db.triples)
	db.triples = append(db.triples, *triple)
	db.stats = append(db.stats, TripleStats{})
	return TripleId(id), nil
}

func (db *inMemoryDatabase) GetTriple(tripleId TripleId) (*Triple, error) {
	if int(tripleId) >= len(db.stats) {
		return nil, &TripleNotFoundError{tripleId}
	}
	return &db.triples[int(tripleId)], nil
}

func (db *inMemoryDatabase) UpdateStats(
	tripleId TripleId, userId UserId,
	stats *TripleStats, vote *VoteStatus, updater Updater) error {
	if int(tripleId) >= len(db.stats) {
		return &TripleNotFoundError{tripleId}
	}

	*stats = db.stats[int(tripleId)]
	*vote = db.vote[tripleUserIdPair{tripleId, userId}]

	update := updater()
	if update {
		db.stats[int(tripleId)] = *stats
		db.vote[tripleUserIdPair{tripleId, userId}] = *vote
	}

	return nil
}

type S struct{}

var _ = Suite(&S{})

func (s *S) TestImageSearchFails(c *C) {
	mfk := MFKImpl{
		Logger: TestLogger{c},
		ImageSearcher: FakeImageSearcher(func(query string) ([]*ImageMetadata, error) {
			c.Check(query, Equals, "query")
			return nil, errors.New("Test error")
		}),
	}

	metadata, err := mfk.ImageSearch("query")
	c.Check(metadata, IsNil)
	c.Check(err, NotNil)
}

func (s *S) TestImageSearchSucceeds(c *C) {
	m1 := ImageMetadata{Url: proto.String("img1"), Context: proto.String("context1")}
	m2 := ImageMetadata{Url: proto.String("img2"), Context: proto.String("context2")}

	mfk := MFKImpl{
		Logger: TestLogger{c},
		ImageSearcher: FakeImageSearcher(func(query string) ([]*ImageMetadata, error) {
			c.Check(query, Equals, "query")
			return []*ImageMetadata{&m1, &m2}, nil
		}),
	}

	metadata, err := mfk.ImageSearch("query")
	c.Check(err, IsNil)
	c.Check(metadata, ProtoEquals, &ImageSearchResponse{
		Image: []*WrappedImageMetadata{
			{Metadata: &m1},
			{Metadata: &m2},
		}})
}

var MAKE_REQUEST = MakeTripleRequest{
	A: &MakeTripleRequest_Entity{
		Name: proto.String("Bulbasaur"),
		Image: &WrappedImageMetadata{
			Metadata: &ImageMetadata{
				Url:     proto.String("bulbapic.png"),
				Context: proto.String("bulbapedia.com/bulbasaur"),
			},
		},
	},
	B: &MakeTripleRequest_Entity{
		Name: proto.String("Pikachu"),
		Image: &WrappedImageMetadata{
			Metadata: &ImageMetadata{
				Url:     proto.String("pikachu.png"),
				Context: proto.String("bulbapedia.com/pikachu"),
			},
		},
	},
	C: &MakeTripleRequest_Entity{
		Name: proto.String("Charmander"),
		Image: &WrappedImageMetadata{
			Metadata: &ImageMetadata{
				Url:     proto.String("char.png"),
				Context: proto.String("bulbapedia.com/charmander"),
			},
		},
	},
}

func (s *S) TestMakeTriple_FetchFails(c *C) {
	mfk := MFKImpl{
		Logger: TestLogger{c},
		ImageFetcher: FakeImageFetcher(func(metadata *ImageMetadata) chan ImageOrError {
			ret := make(chan ImageOrError, 1)
			switch *metadata.Url {
			case "bulbapic.png":
				ret <- ImageOrError{}
			case "pikachu.png":
				ret <- ImageOrError{}
			case "char.png":
				ret <- ImageOrError{error: errors.New("Fetch error")}
			}
			return ret
		}),
	}

	_, err := mfk.MakeTriple(&MAKE_REQUEST)
	c.Check(err, NotNil)
}

func (s *S) TestMakeTriple_Success(c *C) {
	im1 := Image{
		ContentType: proto.String("image/png"),
		Data:        []byte("bulb"),
	}
	im2 := Image{
		ContentType: proto.String("image/png"),
		Data:        []byte("pik"),
	}
	im3 := Image{
		ContentType: proto.String("image/png"),
		Data:        []byte("char"),
	}

	db := NewInMemoryDb()

	mfk := MFKImpl{
		UserId: "Scott",
		Logger: TestLogger{c},
		ImageFetcher: FakeImageFetcher(func(metadata *ImageMetadata) chan ImageOrError {
			ret := make(chan ImageOrError, 1)
			switch *metadata.Url {
			case "bulbapic.png":
				ret <- ImageOrError{Image: &im1}
			case "pikachu.png":
				ret <- ImageOrError{Image: &im2}
			case "char.png":
				ret <- ImageOrError{Image: &im3}
			}
			return ret
		}),
		Database: db,
	}

	resp, err := mfk.MakeTriple(&MAKE_REQUEST)
	c.Check(resp, ProtoEquals, &MakeTripleResponse{
		TripleId: proto.Int64(0),
	})
	c.Check(err, IsNil)

	c.Check(len(db.triples), Equals, 1)
	c.Check(&db.triples[0], ProtoEquals, &Triple{
		CreatorId: proto.String("Scott"),
		A: &Triple_Entity{
			Name:  MAKE_REQUEST.A.Name,
			Image: &im1,
		},
		B: &Triple_Entity{
			Name:  MAKE_REQUEST.B.Name,
			Image: &im2,
		},
		C: &Triple_Entity{
			Name:  MAKE_REQUEST.C.Name,
			Image: &im3,
		},
	})
}

// The ProtoEquals checker verifies that the obtained and expected values are
// both Google protocol buffers, and that they are proto-equal.
var EqualsStats Checker = &equalsStatsChecker{
	&CheckerInfo{Name: "EqualsStats", Params: []string{
			"obtained", "skips",
			"am", "af", "ak",
			"bm", "bf", "bk",
			"cm", "cf", "ck",
	}},
}

type equalsStatsChecker struct {
	*CheckerInfo
}

func (checker *equalsStatsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained, ok := params[0].(TripleStats)
	if !ok {
		return false, "obtained is not a TripleStats"
	}

	i := func(num interface{}) uint64 {
		return uint64(num.(int))
	}

	return obtained.Skips == i(params[1]) &&
		obtained.A.Marry == i(params[2]) &&
		obtained.A.Fuck == i(params[3]) &&
		obtained.A.Kill == i(params[4]) &&
		obtained.B.Marry == i(params[5]) &&
		obtained.B.Fuck == i(params[6]) &&
		obtained.B.Kill == i(params[7]) &&
		obtained.C.Marry == i(params[8]) &&
		obtained.C.Fuck == i(params[9]) &&
		obtained.C.Kill == i(params[10]), "";
}

func makeStats(
	am, af, ak,
	bm, bf, bk,
	cm, cf, ck uint64) TripleStats{
	return TripleStats{
	A: Tally {
			Marry: am,
			Fuck: af,
			Kill: ak,
		},
	B: Tally {
			Marry: bm,
			Fuck: bf,
			Kill: bk,
		},
	C: Tally {
			Marry: cm,
			Fuck: cf,
			Kill: ck,
		},
	}
}

func (s *S) TestBasicVoting(c *C) {
	db := NewInMemoryDb()
	t1, _ := db.AddTriple(&Triple{})
	db.AddTriple(&Triple{})

	mfk1 := MFKImpl{
		UserId: "Scott",
		Logger: TestLogger{c},
		Database: db,
	}
	mfk2 := MFKImpl{
		UserId: "Mike",
		Logger: TestLogger{c},
		Database: db,
	}

	// All clear
	stats, vote, err := mfk1.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0)
	c.Check(vote, Equals, VoteStatus_UNSET)

	stats, vote, err = mfk2.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0)
	c.Check(vote, Equals, VoteStatus_UNSET)

	// Change Scott's vote
	err = mfk1.ChangeVote(t1, VoteStatus_MFK)
	c.Check(err, IsNil)

	stats, vote, err = mfk1.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0)
	c.Check(vote, Equals, VoteStatus_MFK)

	stats, vote, err = mfk2.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1)
	c.Check(vote, Equals, VoteStatus_UNSET)

	// Mike skips the first time.
	err = mfk2.ChangeVote(t1, VoteStatus_SKIP)
	c.Check(err, IsNil)

	stats, vote, err = mfk1.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0)
	c.Check(vote, Equals, VoteStatus_MFK)

	stats, vote, err = mfk2.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1)
	c.Check(vote, Equals, VoteStatus_SKIP)

	// Then votes.
	err = mfk2.ChangeVote(t1, VoteStatus_MKF)
	c.Check(err, IsNil)

	stats, vote, err = mfk1.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		1, 0, 0,
		0, 0, 1,
		0, 1, 0)
	c.Check(vote, Equals, VoteStatus_MFK)

	stats, vote, err = mfk2.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1)
	c.Check(vote, Equals, VoteStatus_MKF)

	// The stats from some 3rd player's perspective.
	mfk3 := MFKImpl{
		UserId: "Skwisgaar",
		Logger: TestLogger{c},
		Database: db,
	}

	stats, vote, err = mfk3.GetTripleStatsForUser(t1)
	c.Check(err, IsNil)
	c.Check(stats, EqualsStats, 0,
		2, 0, 0,
		0, 1, 1,
		0, 1, 1)
	c.Check(vote, Equals, VoteStatus_UNSET)
}

func (s *S) TestVoteWhenNoTriple(c *C) {
	db := NewInMemoryDb()

	mfk := MFKImpl{
		UserId: "Scott",
		Logger: TestLogger{c},
		Database: db,
	}

	err := mfk.ChangeVote(TripleId(3), VoteStatus_MFK)
	c.Check(*(err.(*TripleNotFoundError)), Equals, TripleNotFoundError{
		TripleId: 3,
	})
}

func (s *S) TestGetStatsNoTriple(c *C) {
	db := NewInMemoryDb()

	mfk := MFKImpl{
		UserId: "Scott",
		Logger: TestLogger{c},
		Database: db,
	}

	_, _, err := mfk.GetTripleStatsForUser(TripleId(3))
	c.Check(*(err.(*TripleNotFoundError)), Equals, TripleNotFoundError{
		TripleId: 3,
	})
}