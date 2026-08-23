package repository

import (
	"fmt"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindGet
	ErrorKindNotFound
	ErrorKindCreate
	ErrorKindUpdate
	ErrorKindDelete

	errorKindUnknownString  = "unknown"
	errorKindGetString      = "get"
	errorKindNotFoundString = "not_found"
	errorKindCreateString   = "create"
	errorKindUpdateString   = "update"
	errorKindDeleteString   = "delete"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindGet:
		return errorKindGetString
	case ErrorKindNotFound:
		return errorKindNotFoundString
	case ErrorKindCreate:
		return errorKindCreateString
	case ErrorKindUpdate:
		return errorKindUpdateString
	case ErrorKindDelete:
		return errorKindDeleteString
	default:
		return errorKindUnknownString
	}
}

type Error struct {
	Kind  ErrorKind
	error error
}

func newError(err error, kind ErrorKind) Error {
	return Error{error: err, Kind: kind}
}

func NewUnknownError(err error) Error {
	return newError(err, ErrorKindUnknown)
}

func NewGetError(err error) Error {
	return newError(err, ErrorKindGet)
}

func NewNotFoundError(err error) Error {
	return newError(err, ErrorKindNotFound)
}

func NewCreateError(err error) Error {
	return newError(err, ErrorKindCreate)
}

func NewUpdateError(err error) Error {
	return newError(err, ErrorKindUpdate)
}

func NewDeleteError(err error) Error {
	return newError(err, ErrorKindDelete)
}

func (e Error) Error() string {
	return fmt.Sprintf("customer repository error (kind = %s): %v", e.Kind.String(), e.error)
}
