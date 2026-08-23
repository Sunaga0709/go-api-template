package commandservice

import (
	"errors"
	"fmt"

	bookuc "github.com/Sunaga0709/go-api-template/internal/domain/book/usecase"
	cfbuc "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/usecase"
	customeruc "github.com/Sunaga0709/go-api-template/internal/domain/customer/usecase"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindInvalidValue
	ErrorKindNotFound
	ErrorKindExecUsecase
	ErrorKindAlreadyExists

	errorKindUnknownString       = "unknown"
	errorKindInvalidValueString  = "invalid_value"
	errorKindNotFoundString      = "not_found"
	errorKindExecUsecaseString   = "exec_usecase"
	errorKindAlreadyExistsString = "already_exists"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindInvalidValue:
		return errorKindInvalidValueString
	case ErrorKindNotFound:
		return errorKindNotFoundString
	case ErrorKindExecUsecase:
		return errorKindExecUsecaseString
	case ErrorKindAlreadyExists:
		return errorKindAlreadyExistsString
	default:
		return errorKindUnknownString
	}
}

type Error struct {
	Kind          ErrorKind
	error         error
	ClientMessage *string
}

func newError(err error, kind ErrorKind) Error {
	return Error{
		Kind:          kind,
		error:         err,
		ClientMessage: nil,
	}
}

func newUnknownError(err error) Error {
	return newError(err, ErrorKindUnknown)
}

func newErrorFromCustomerUsecaseError(err error) Error {
	if ucErr, ok := errors.AsType[customeruc.Error](err); ok {
		clientMessage := ucErr.ClientMessage
		var errKind ErrorKind
		switch ucErr.Kind {
		case customeruc.ErrorKindInvalidValue:
			errKind = ErrorKindInvalidValue
		case customeruc.ErrorKindNotFound:
			errKind = ErrorKindNotFound
		default:
			errKind = ErrorKindExecUsecase
		}

		e := newError(err, errKind)
		if clientMessage != nil {
			return e.WithClientMessage(*clientMessage)
		}

		return e
	}

	return newUnknownError(err)
}

func newErrorFromCustomerFavoriteBookUsecaseError(err error) Error {
	if ucErr, ok := errors.AsType[cfbuc.Error](err); ok {
		clientMessage := ucErr.ClientMessage
		var errKind ErrorKind
		switch ucErr.Kind {
		case cfbuc.ErrorKindInvalidValue:
			errKind = ErrorKindInvalidValue
		case cfbuc.ErrorKindAlreadyFavorited:
			errKind = ErrorKindAlreadyExists
		case cfbuc.ErrorKindNotFavorited:
			errKind = ErrorKindNotFound
		default:
			errKind = ErrorKindExecUsecase
		}

		e := newError(err, errKind)
		if clientMessage != nil {
			return e.WithClientMessage(*clientMessage)
		}

		return e
	}

	return newUnknownError(err)
}

func newErrorFromBookUsecaseError(err error) Error {
	if ucErr, ok := errors.AsType[bookuc.Error](err); ok {
		clientMessage := ucErr.ClientMessage
		var errKind ErrorKind
		switch ucErr.Kind {
		case bookuc.ErrorKindNotFound:
			errKind = ErrorKindNotFound
		case bookuc.ErrorKindInvalidValue:
			errKind = ErrorKindInvalidValue
		default:
			errKind = ErrorKindExecUsecase
		}

		e := newError(err, errKind)
		if clientMessage != nil {
			return e.WithClientMessage(*clientMessage)
		}

		return e
	}

	return newUnknownError(err)
}

func (e Error) WithClientMessage(msg string) Error {
	err := e
	err.ClientMessage = &msg
	return err
}

func ParseError(err error) (Error, bool) {
	return errors.AsType[Error](err)
}

func (e Error) Error() string {
	return fmt.Sprintf("commandservice error (kind=%s): %v", e.Kind.String(), e.error)
}
