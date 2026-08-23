package model

import (
	"fmt"

	"github.com/google/uuid"
)

type BookID struct {
	string
}

func GenerateBookID() (BookID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return BookID{}, newUnknownError(fmt.Errorf("failed to generate book id: %w", err))
	}

	return BookID{id.String()}, nil
}

func NewBookID(id string) (BookID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return BookID{}, newInvalidValueError(fmt.Errorf("invalid book id: %w", err))
	}

	return BookID{id}, nil
}

func (b BookID) String() string {
	return b.string
}
