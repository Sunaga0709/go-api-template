package model

import "fmt"

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindInvalidValue

	errorKindUnknown      = "unknown"
	errorKindInvalidValue = "invalid_value"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindInvalidValue:
		return errorKindInvalidValue
	default:
		return errorKindUnknown
	}
}

type Error struct {
	error error
	Kind  ErrorKind
}

func newError(err error, kind ErrorKind) Error {
	return Error{
		error: err,
		Kind:  kind,
	}
}

func newInvalidValueError(err error) Error {
	return newError(err, ErrorKindInvalidValue)
}

func newUnknownError(err error) Error {
	return newError(err, ErrorKindUnknown)
}

func (e Error) Error() string {
	return fmt.Sprintf("customer model error (kind=%s): %v", e.Kind.String(), e.error)
}
