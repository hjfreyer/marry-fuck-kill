package mfklib

import (
	"fmt"
)

type TripleId int64
type UserId string

type EntityNotFoundError struct {
	Type string
	Id   int64
	Err  error
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("Entity of type %q with id %d not found. Error: %s",
		e.Type, e.Id, e.Err)
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

type Updater func() (bool, error)

type Database interface {
	AddTriple(*Triple) (TripleId, error)
	GetTriple(tripleId TripleId) (*Triple, error)

	UpdateStats(TripleId, UserId, *TripleStats, *TripleUserStatus, Updater) error
}

