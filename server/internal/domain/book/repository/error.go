package repository

import "fmt"

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindGet
	ErrorKindNotFound
	ErrorKindCreate

	errorKindUnknownString  = "unknown"
	errorKindGetString      = "get"
	errorKindNotFoundString = "not_found"
	errorKindCreateString   = "create"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindGet:
		return errorKindGetString
	case ErrorKindNotFound:
		return errorKindNotFoundString
	case ErrorKindCreate:
		return errorKindCreateString
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

func NewNotFoundError(err error) Error {
	return newError(ErrorKindNotFound, err)
}

func NewCreateError(err error) Error {
	return newError(ErrorKindCreate, err)
}

func NewUnknownError(err error) Error {
	return newError(ErrorKindUnknown, err)
}

func (e Error) Error() string {
	return fmt.Sprintf("book repository error (kind = %s): %v", e.Kind.String(), e.error)
}
