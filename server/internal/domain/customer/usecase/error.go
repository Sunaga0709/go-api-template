package usecase

import (
	"errors"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/domain/customer/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/repository"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindInvalidValue
	ErrorKindNotFound
	ErrorKindExecRepository

	errorKindUnknownString        = "unknown"
	errorKindInvalidValueString   = "invalid_value"
	errorKindNotFoundString       = "not_found"
	errorKindExecRepositoryString = "exec_repository"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindInvalidValue:
		return errorKindInvalidValueString
	case ErrorKindNotFound:
		return errorKindNotFoundString
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

func _newError(err error, kind ErrorKind) Error {
	return Error{
		Kind:          kind,
		error:         err,
		ClientMessage: nil,
	}
}

//nolint:unused // allow unused
func newError(err error) Error {
	if modelErr, ok := errors.AsType[model.Error](err); ok {
		return newErrorFromModelError(modelErr)
	}

	if repoErr, ok := errors.AsType[repository.Error](err); ok {
		return newErrorFromRepositoryError(repoErr)
	}

	return newUnknownError(err)
}

func newInvalidValueError(err error) Error {
	return _newError(err, ErrorKindInvalidValue)
}

func newNotFoundError(err error) Error {
	return _newError(err, ErrorKindNotFound)
}

func newExecRepositoryError(err error) Error {
	return _newError(err, ErrorKindExecRepository)
}

func newUnknownError(err error) Error {
	return _newError(err, ErrorKindUnknown)
}

func newErrorFromModelError(err error) Error {
	if modelErr, ok := errors.AsType[model.Error](err); ok {
		switch modelErr.Kind {
		case model.ErrorKindInvalidValue:
			return newInvalidValueError(err)
		default:
			return newUnknownError(err)
		}
	}

	return newUnknownError(err)
}

func newErrorFromRepositoryError(err error) Error {
	if repoErr, ok := errors.AsType[repository.Error](err); ok {
		switch repoErr.Kind {
		case repository.ErrorKindNotFound:
			return newNotFoundError(err)
		case repository.ErrorKindGet, repository.ErrorKindCreate, repository.ErrorKindUpdate, repository.ErrorKindDelete:
			return newExecRepositoryError(err)
		default:
			return newUnknownError(err)
		}
	}

	return newUnknownError(err)
}

func ParseError(err error) (Error, bool) {
	return errors.AsType[Error](err)
}

func (e Error) WithClientMessage(msg string) Error {
	err := e
	err.ClientMessage = &msg
	return err
}

func (e Error) Error() string {
	return fmt.Sprintf("customer usecase error (kind=%s): %v", e.Kind.String(), e.error)
}
