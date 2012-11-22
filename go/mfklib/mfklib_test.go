package mfklib

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestLogger struct{ *C }

func (t TestLogger) Infof(format string, args ...interface{}) {
	t.Logf(format, args...)
}

func (t TestLogger) Warningf(format string, args ...interface{}) {
	t.Logf("W: "+format, args...)
}

type FakeImageSearcher func(query string) (results []*ImageMetadata, err error)

func (f FakeImageSearcher) Search(query string) (results []*ImageMetadata, err error) {
	return f(query)
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
