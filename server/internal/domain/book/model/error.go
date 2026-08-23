package model

import "fmt"

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindInvalidValue
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindInvalidValue:
		return "invalid_value"
	default:
		return "unknown"
	}
}

type Error struct {
	Kind  ErrorKind
	error error
}

func newError(err error, kind ErrorKind) Error {
	return Error{
		Kind:  kind,
		error: err,
	}
}

func newInvalidValueError(err error) Error {
	return newError(err, ErrorKindInvalidValue)
}

func newUnknownError(err error) Error {
	return newError(err, ErrorKindUnknown)
}

func (e Error) Error() string {
	return fmt.Sprintf("book model error (kind = %s): %v", e.Kind.String(), e.error)
}
