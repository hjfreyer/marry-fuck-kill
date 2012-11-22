package mfklib

import (
	. "launchpad.net/gocheck"
	"testing"
	// "fmt"
	"code.google.com/p/goprotobuf/proto"
	"errors"
)

func Test(t *testing.T) { TestingT(t) }

type TestLogger struct{ *C }

func (t TestLogger) Infof(format string, args ...interface{}) {
	t.Logf(format, args...)
}

func (t TestLogger) Warningf(format string, args ...interface{}) {
	t.Logf("W: "+format, args...)
}

type FakeImageSearcher func(query string) (results []ImageMetadata, err error)

func (f FakeImageSearcher) Search(query string) (results []ImageMetadata, err error) {
	return f(query)
}

// type ImageSearchImpl struct {
// 	t *testing.T
// 	fails, count int
// }

// func (i *ImageSearchImpl) Search(query string)  (results []ImageMetadata, err error) {
// 	if query != "query" {
// 		i.t.Errorf("Illegal query: %s", query)
// 	}

// 	if i.count < i.fails {
// 		err = fmt.Errorf("Test error: %d", i.count)
// 	} else {
// 		results = []ImageMetadata{
// 			{Url: proto.String("img1"), Context: proto.String("context1")},
// 			{Url: proto.String("img2"), Context: proto.String("context2")},
// 		}
// 	}
// 	i.count++
// 	return
// }

type S struct{}

var _ = Suite(&S{})

// func (s *S) TestHelloWorld(c *C) {
//     c.Check(42, Equals, "42")
// //    c.Check(os.Errno(13), Matches, "perm.*accepted")
// }

func assertProtoEquals(t *testing.T, expected, actual proto.Message) {
	if !proto.Equal(expected, actual) {
		t.Errorf("Protos not equal, expected:\n%s\nActual:\n%s",
			proto.MarshalTextString(expected),
			proto.MarshalTextString(actual))
	}
}

func (s *S) TestImageSearchFails(c *C) {
	err := errors.New("Test error")

	mfk := MFKImpl{
		Logger: TestLogger{c},
		ImageSearcher: FakeImageSearcher(func(query string) ([]ImageMetadata, error) {
			return nil, err
		}),
	}

	metadata, actual_error := mfk.ImageSearch("query")
	c.Check(metadata, IsNil)
	c.Check(actual_error, Equals, err)

	// //	if metadata !=
	// 	if err == nil {
	// 		t.Error("expected error")
	// 	}

	// 	_, err := mfk.ImageSearch("query", 1)

}
