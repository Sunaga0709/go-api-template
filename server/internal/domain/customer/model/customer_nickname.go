package model

import (
	"fmt"
	"unicode/utf8"
)

const (
	minNicknameLength = 1
	maxNicknameLength = 50

	defaultNickname = "未設定"
)

type CustomerNickname struct {
	string
}

func NewCustomerNickname(nickname string) (CustomerNickname, error) {
	cn := CustomerNickname{nickname}
	if err := cn.validate(); err != nil {
		return CustomerNickname{}, err
	}

	return cn, nil
}

func DefaultCustomerNickname() CustomerNickname {
	return CustomerNickname{defaultNickname}
}

func (c *CustomerNickname) validate() error {
	length := utf8.RuneCountInString(c.string)
	if length < minNicknameLength || maxNicknameLength < length {
		return newInvalidValueError(fmt.Errorf("invalid customer nickname: got = %s", c.string))
	}

	return nil
}

func (c *CustomerNickname) String() string {
	return c.string
}
