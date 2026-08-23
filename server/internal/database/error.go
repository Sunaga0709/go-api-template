package database

import "fmt"

type Error struct {
	error error
}

func newError(err error) Error {
	return Error{
		error: err,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("database error: %v", e.error)
}
