package model

import (
	"fmt"
	"time"

	"github.com/Sunaga0709/go-api-template/internal/location"
)

type CustomerLocation struct {
	value *time.Location
}

// NewCustomerLocation `nil`が入力された場合は、デフォルトのロケーション(`JST`)となる
func NewCustomerLocation(loc *time.Location) CustomerLocation {
	if loc == nil {
		return DefaultCustomerLocation()
	}

	return CustomerLocation{value: loc}
}

func DefaultCustomerLocation() CustomerLocation {
	return CustomerLocation{
		value: location.JST(),
	}
}

func NewCustomerLocationFromString(loc string) (CustomerLocation, error) {
	location, err := time.LoadLocation(loc)
	if err != nil {
		return CustomerLocation{}, newInvalidValueError(fmt.Errorf("invalid customer location: got = %s", loc))
	}

	return CustomerLocation{value: location}, nil
}

// Value `nil`を保持していた場合は`JST`を返却する
func (c *CustomerLocation) Location() *time.Location {
	if c.value == nil {
		return DefaultCustomerLocation().value
	}

	return c.value
}

func (c *CustomerLocation) String() string {
	if c.value == nil {
		return DefaultCustomerLocation().value.String()
	}

	return c.value.String()
}
