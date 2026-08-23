package model

import (
	"fmt"
	"time"

	"github.com/Sunaga0709/go-api-template/internal/date"
)

// CustomerBirthday `time.Time`型を使用するが時分秒、ロケーションは扱わない
// ロケーションは`UTC`に統一して保持する
type CustomerBirthday struct {
	value date.Date
}

func NewCustomerBirthday(birthday date.Date) CustomerBirthday {
	return CustomerBirthday{value: birthday}
}

func NewCustomerBirthdayFromStdTime(birthday time.Time) (CustomerBirthday, error) {
	cbd, err := date.NewDate(birthday.Year(), birthday.Month(), birthday.Day())
	if err != nil {
		return CustomerBirthday{}, newInvalidValueError(fmt.Errorf("invalid customer birthday: got = %v", birthday))
	}

	return CustomerBirthday{cbd}, nil
}

func DefaultCustomerBirthday() CustomerBirthday {
	return CustomerBirthday{date.Min()}
}

func (c *CustomerBirthday) IsDefault() bool {
	return c.value.IsMin()
}

func (c *CustomerBirthday) IsPast() bool {
	return c.value.IsPast(time.Now())
}

func (c *CustomerBirthday) Date() date.Date {
	return c.value
}

func (c *CustomerBirthday) ToStdTime() time.Time {
	return time.Date(c.value.Year(), c.value.Month(), c.value.Day(), 0, 0, 0, 0, time.UTC)
}
