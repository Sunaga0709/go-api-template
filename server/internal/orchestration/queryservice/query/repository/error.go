package repository

import "fmt"

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindGet
	ErrorKindNotFound
	ErrorKindConvertDomainModel

	errorKindUnknownString          = "unknown"
	errorKindGetString              = "get"
	errorKindNotFoundString         = "not_found"
	errKindConvertDomainModelString = "convert_domain_model"
)

func (e ErrorKind) String() string {
	switch e {
	case ErrorKindGet:
		return errorKindGetString
	case ErrorKindNotFound:
		return errorKindNotFoundString
	case ErrorKindConvertDomainModel:
		return errKindConvertDomainModelString
	default:
		return errorKindUnknownString
	}
}

type Error struct {
	error error
	Kind  ErrorKind
}

func NewError(err error, kind ErrorKind) Error {
	return Error{
		error: err,
		Kind:  kind,
	}
}

func NewUnknownError(err error) Error {
	return NewError(err, ErrorKindUnknown)
}

func NewGetError(err error) Error {
	return NewError(err, ErrorKindGet)
}

func NewNotFoundError(err error) Error {
	return NewError(err, ErrorKindNotFound)
}

func NewConvertDomainModelError(err error) Error {
	return NewError(err, ErrorKindConvertDomainModel)
}

func (e Error) Error() string {
	return fmt.Sprintf("query error (kind=%s): %v", e.Kind.String(), e.error)
}
