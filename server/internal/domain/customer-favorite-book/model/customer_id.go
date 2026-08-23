package model

import "errors"

type CustomerID struct {
	string
}

func NewCustomerID(customerID string) (CustomerID, error) {
	id := CustomerID{customerID}
	if err := id.validate(); err != nil {
		return CustomerID{}, err
	}

	return id, nil
}

func (c *CustomerID) validate() error {
	if c.string == "" {
		return newInvalidValueError(errors.New("customer id is empty"))
	}

	return nil
}

func (c *CustomerID) String() string {
	return c.string
}
