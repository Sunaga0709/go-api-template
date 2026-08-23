package model

import (
	"time"
)

type CustomerRegisteredAt struct {
	value time.Time
}

func NewCustomerRegisteredAt(registeredAt time.Time) CustomerRegisteredAt {
	return CustomerRegisteredAt{value: registeredAt.UTC()}
}

func DefaultCustomerRegisteredAt() CustomerRegisteredAt {
	return CustomerRegisteredAt{time.Time{}}
}

func (c *CustomerRegisteredAt) Time() time.Time {
	return c.value
}

func (c *CustomerRegisteredAt) IsRegistered() bool {
	return !c.value.IsZero()
}

func (c *CustomerRegisteredAt) Register() {
	c.value = time.Now().UTC()
}
