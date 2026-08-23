package model

import (
	"errors"
	"time"
)

type Customers []*Customer

type Customer struct {
	CustomerID   CustomerID
	Nickname     CustomerNickname
	Birthday     CustomerBirthday
	Location     CustomerLocation
	RegisteredAt CustomerRegisteredAt
}

func NewCustomer(
	customerID CustomerID,
	nickname CustomerNickname,
	birthday CustomerBirthday,
	location CustomerLocation,
	registeredAt CustomerRegisteredAt,
) (*Customer, error) {
	customer := Customer{
		CustomerID:   customerID,
		Nickname:     nickname,
		Birthday:     birthday,
		Location:     location,
		RegisteredAt: registeredAt,
	}
	if err := customer.validate(); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (c *Customer) validate() error {
	if c.RegisteredAt.IsRegistered() && c.Birthday.IsDefault() && !c.Birthday.IsPast() {
		return newInvalidValueError(errors.New("not set customer birthday when registered"))
	}

	return nil
}

// Update 顧客情報の更新を行う
// `RegisteredAt`がゼロ値（プロフィール未登録）の場合、現在日時を設定する
func (c *Customer) Update(nickname CustomerNickname, birthday CustomerBirthday, location CustomerLocation) error {
	registeredAt := c.RegisteredAt
	if !registeredAt.IsRegistered() {
		registeredAt = NewCustomerRegisteredAt(time.Now().UTC())
	}

	updated, err := NewCustomer(c.CustomerID, nickname, birthday, location, registeredAt)
	if err != nil {
		return err
	}

	c.Nickname = updated.Nickname
	c.Birthday = updated.Birthday
	c.Location = updated.Location
	c.RegisteredAt = updated.RegisteredAt

	return nil
}

func (c *Customer) IsRegistered() bool {
	return c.RegisteredAt.IsRegistered()
}
