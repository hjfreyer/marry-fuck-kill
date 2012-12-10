package mfklib

import (
	"errors"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
	. "launchpad.net/gocheck"
	"testing"
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
type FakeImageFetcher func(url string) chan ImageOrError

func (f FakeImageSearcher) Search(query string) (results []*ImageMetadata, err error) {
	return f(query)
}

func (f FakeImageFetcher) FetchImage(url string) chan ImageOrError {
	return f(url)
}

type inMemoryDatabase struct {
	triples []*Triple
	stats   []*TripleStats
}

func (db *inMemoryDatabase) AddTriple(triple *Triple) (int64, error) {
	id := len(db.triples)
	db.triples = append(db.triples, triple)
	db.stats = append(db.stats, &TripleStats{})
	return int64(id), nil
}

func (db *inMemoryDatabase) GetTriple(id int64) (*Triple, error) {
	return db.triples[id], nil
}

func (db *inMemoryDatabase) UpdateStats(id int64, updater TripleStatsUpdater) error {
	updater(db.stats[id])
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
		ImageFetcher: FakeImageFetcher(func(url string) chan ImageOrError {
			ret := make(chan ImageOrError, 1)
			switch url {
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

	db := inMemoryDatabase{}

	mfk := MFKImpl{
		UserId: "Scott",
		Logger: TestLogger{c},
		ImageFetcher: FakeImageFetcher(func(url string) chan ImageOrError {
			ret := make(chan ImageOrError, 1)
			switch url {
			case "bulbapic.png":
				ret <- ImageOrError{Image: &im1}
			case "pikachu.png":
				ret <- ImageOrError{Image: &im2}
			case "char.png":
				ret <- ImageOrError{Image: &im3}
			}
			return ret
		}),
		Database: &db,
	}

	resp, err := mfk.MakeTriple(&MAKE_REQUEST)
	c.Check(resp, ProtoEquals, &MakeTripleResponse{
		TripleId: proto.Int64(0),
	})
	c.Check(err, IsNil)

	c.Check(len(db.triples), Equals, 1)
	c.Check(db.triples[0], ProtoEquals, &Triple{
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
