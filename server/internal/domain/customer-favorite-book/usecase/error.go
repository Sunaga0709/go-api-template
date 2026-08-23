package usecase

import (
	"errors"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/repository"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindInvalidValue
	ErrorKindExecRepository
	ErrorKindAlreadyFavorited
	ErrorKindNotFavorited

	errorKindUnknownString          = "unknown"
	errorKindInvalidValueString     = "invalid_value"
	errorKindExecRepositoryString   = "exec_repository"
	errorKindAlreadyFavoritedString = "already_favorited"
	errorKindNotFavoritedString     = "not_favorited"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindInvalidValue:
		return errorKindInvalidValueString
	case ErrorKindExecRepository:
		return errorKindExecRepositoryString
	case ErrorKindAlreadyFavorited:
		return errorKindAlreadyFavoritedString
	case ErrorKindNotFavorited:
		return errorKindNotFavoritedString
	default:
		return errorKindUnknownString
	}
}

type Error struct {
	Kind          ErrorKind
	error         error
	ClientMessage *string
}

func newError(err error) Error {
	if modelErr, ok := errors.AsType[model.Error](err); ok {
		return newErrorFromModelError(modelErr)
	}

	if repoErr, ok := errors.AsType[repository.Error](err); ok {
		return newErrorFromRepositoryError(repoErr)
	}

	return newUnknownError(err)
}

func _newError(kind ErrorKind, err error) Error {
	return Error{
		Kind:          kind,
		error:         err,
		ClientMessage: nil,
	}
}

func newInvalidValueError(err error) Error {
	return _newError(ErrorKindInvalidValue, err)
}

func newExecRepositoryError(err error) Error {
	return _newError(ErrorKindExecRepository, err)
}

func newAlreadyFavoritedError(err error) Error {
	return _newError(ErrorKindAlreadyFavorited, err)
}

func newNotFavoritedError(err error) Error {
	return _newError(ErrorKindNotFavorited, err)
}

func newUnknownError(err error) Error {
	return _newError(ErrorKindUnknown, err)
}

func newErrorFromModelError(err error) Error {
	modelErr, ok := errors.AsType[model.Error](err)
	if !ok {
		return newUnknownError(err)
	}

	switch modelErr.Kind {
	case model.ErrorKindInvalidValue:
		return newInvalidValueError(err)
	case model.ErrorKindAlreadyFavorited:
		return newAlreadyFavoritedError(err)
	case model.ErrorKindNotFavorited:
		return newNotFavoritedError(err)
	default:
		return newUnknownError(err)
	}
}

func newErrorFromRepositoryError(err error) Error {
	repoErr, ok := errors.AsType[repository.Error](err)
	if !ok {
		return newUnknownError(err)
	}

	switch repoErr.Kind {
	case repository.ErrorKindUnknown:
		return newUnknownError(err)
	default:
		return newExecRepositoryError(err)
	}
}

func ParseError(err error) (Error, bool) {
	return errors.AsType[Error](err)
}

func (e Error) WithClientMessage(msg string) Error {
	return Error{
		Kind:          e.Kind,
		error:         e.error,
		ClientMessage: &msg,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("customer favorite book usecase error (kind = %s): %v", e.Kind.String(), e.error)
}
