package usecase

import (
	"errors"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/domain/book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/repository"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindNotFound
	ErrorKindInvalidValue
	ErrorKindExecRepository

	errorKindUnknownString        = "unknown"
	errorKindNotFoundString       = "not_found"
	errorKindInvalidValueString   = "invalid_value"
	errorKindExecRepositoryString = "exec_repository"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindNotFound:
		return errorKindNotFoundString
	case ErrorKindInvalidValue:
		return errorKindInvalidValueString
	case ErrorKindExecRepository:
		return errorKindExecRepositoryString
	default:
		return errorKindUnknownString
	}
}

type Error struct {
	Kind          ErrorKind
	error         error
	ClientMessage *string
}

func _newError(kind ErrorKind, err error) Error {
	return Error{
		Kind:          kind,
		error:         err,
		ClientMessage: nil,
	}
}

func newErrorFromModelError(err error) Error {
	modelErr, ok := errors.AsType[model.Error](err)
	if !ok {
		return _newError(ErrorKindUnknown, err)
	}

	switch modelErr.Kind {
	case model.ErrorKindInvalidValue:
		return _newError(ErrorKindInvalidValue, err)
	default:
		return _newError(ErrorKindUnknown, err)
	}
}

func newErrorFromRepositoryError(err error) Error {
	repoErr, ok := errors.AsType[repository.Error](err)
	if !ok {
		return _newError(ErrorKindUnknown, err)
	}

	switch repoErr.Kind {
	case repository.ErrorKindNotFound:
		return _newError(ErrorKindNotFound, err)
	default:
		return _newError(ErrorKindExecRepository, err)
	}
}

func ParseError(err error) (Error, bool) {
	return errors.AsType[Error](err)
}

func (e Error) Error() string {
	return fmt.Sprintf("book usecase error (kind = %s): %v", e.Kind.String(), e.error)
}

func (e Error) WithClientMessage(msg string) Error {
	return Error{
		Kind:          e.Kind,
		error:         e.error,
		ClientMessage: &msg,
	}
}
