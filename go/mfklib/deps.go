package mfklib

import (
	"fmt"
_	"errors"
)

type TripleId int64
type UserId string

type TripleNotFoundError struct {
	TripleId
}

func (e TripleNotFoundError) Error() string {
	return fmt.Sprintf("Triple not found: %d", e.TripleId)
}

type LowerLevelError struct {
	Message string
	Err error
}

func (e LowerLevelError) Error() string {
	return fmt.Sprintf("Error while performing task %q: %s", e.Message, e.Err)
}

func NewLowerLevelError(err error, format string, a ...interface{}) *LowerLevelError {
	return &LowerLevelError{fmt.Sprintf(format, a...), err}
}

type IllegalArgumentError struct {
	Func     string
	Argument string
	Value    interface{}
	Cause    string
}

func (e IllegalArgumentError) Error() string {
	return fmt.Sprintf("Function %q got illegal value %q for argument %q: %s",
		e.Func, e.Value, e.Argument, e.Cause)
}

type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
}

// ImageSearcher is an interface for a backend component that searches for an
// image on Google Image Search, given a query string. If there is an error, It
// will be of type *ImageSearchError.
type ImageSearcher interface {
	Search(query string) ([]*ImageMetadata, error)
}

type ImageSearchError struct {
	Query string
	Err error
}

func (e ImageSearchError) Error() string {
	return fmt.Sprintf("Error while looking for images for query %q: %s", e.Query, e.Err)
}

type ImageFetcher interface {
	FetchImage(metadata *ImageMetadata) chan ImageOrError
}

type ImageOrError struct {
	*Image
	error
}

type Updater func() bool

type Database interface {
	AddTriple(*Triple) (TripleId, error)
	GetTriple(tripleId TripleId) (*Triple, error)

	UpdateStats(TripleId, UserId, *TripleStats, *VoteStatus, Updater) error
}
