package repository

import "fmt"

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindGet
	ErrorKindUpdate

	errorKindUnknownString = "unknown"
	errorKindGetString     = "get"
	errorKindUpdateString  = "update"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindGet:
		return errorKindGetString
	case ErrorKindUpdate:
		return errorKindUpdateString
	default:
		return errorKindUnknownString
	}
}

type Error struct {
	Kind  ErrorKind
	error error
}

func newError(kind ErrorKind, err error) Error {
	return Error{
		Kind:  kind,
		error: err,
	}
}

func NewGetError(err error) Error {
	return newError(ErrorKindGet, err)
}

func NewUpdateError(err error) Error {
	return newError(ErrorKindUpdate, err)
}

func NewUnknownError(err error) Error {
	return newError(ErrorKindUnknown, err)
}

func (e Error) Error() string {
	return fmt.Sprintf("customer favorite book repository error (kind = %s): %v", e.Kind.String(), e.error)
}
