package queryservice

import (
	"errors"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/repository"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindExecQuery
	ErrorKindInvalidInput
	ErrorKindNotFound

	errorKindUnknownString      = "unknown"
	errorKindExecQueryString    = "exec_query"
	errorKindInvalidInputString = "invalid_input"
	errorKindNotFoundString     = "not_found"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindExecQuery:
		return errorKindExecQueryString
	case ErrorKindInvalidInput:
		return errorKindInvalidInputString
	case ErrorKindNotFound:
		return errorKindNotFoundString
	default:
		return errorKindUnknownString
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

func newExecQueryError(err error) Error {
	return newError(err, ErrorKindExecQuery)
}

func newNotFoundError(err error) Error {
	return newError(err, ErrorKindNotFound)
}

//nolint:unused // allow unused
func newInvalidInputError(err error) Error {
	return newError(err, ErrorKindInvalidInput)
}

func newUnknownError(err error) Error {
	return newError(err, ErrorKindUnknown)
}

func newErrorFromQueryRepositoryError(err error) Error {
	if queryErr, ok := errors.AsType[repository.Error](err); ok {
		switch queryErr.Kind {
		case repository.ErrorKindNotFound:
			return newNotFoundError(err)
		case repository.ErrorKindGet:
			return newExecQueryError(err)
		default:
			return newUnknownError(err)
		}
	}

	return newUnknownError(err)
}

func (e Error) Error() string {
	return fmt.Sprintf("queryservice error (kind=%s): %v", e.Kind.String(), e.error)
}
