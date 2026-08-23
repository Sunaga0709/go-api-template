package model

import (
	"fmt"

	"github.com/google/uuid"
)

type CustomerID struct {
	string
}

func GenerateCustomerID() (CustomerID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return CustomerID{}, newUnknownError(fmt.Errorf("failed to generate customer id: %w", err))
	}

	return CustomerID{id.String()}, nil
}

func NewCustomerID(id string) (CustomerID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return CustomerID{}, newInvalidValueError(fmt.Errorf("invalid customer id: %w", err))
	}

	return CustomerID{id}, nil
}

func (c CustomerID) String() string {
	return c.string
}
