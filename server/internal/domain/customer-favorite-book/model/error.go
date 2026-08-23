package model

import "fmt"

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindInvalidValue
	ErrorKindAlreadyFavorited
	ErrorKindNotFavorited

	errorKindUnknownString          = "unknown"
	errorKindInvalidValueString     = "invalid_value"
	errorKindAlreadyFavoritedString = "already_favorited"
	errorKindNotFavoritedString     = "not_favorited"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindInvalidValue:
		return errorKindInvalidValueString
	case ErrorKindAlreadyFavorited:
		return errorKindAlreadyFavoritedString
	case ErrorKindNotFavorited:
		return errorKindNotFavoritedString
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

func newInvalidValueError(err error) Error {
	return newError(ErrorKindInvalidValue, err)
}

func newAlreadyFavoritedError(err error) Error {
	return newError(ErrorKindAlreadyFavorited, err)
}

func newNotFavoritedError(err error) Error {
	return newError(ErrorKindNotFavorited, err)
}

func (e Error) Error() string {
	return fmt.Sprintf("customer favorite book model error (kind = %s): %v", e.Kind.String(), e.error)
}
