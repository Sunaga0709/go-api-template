package model

import (
	"fmt"
	"unicode/utf8"
)

const (
	minTitleLength = 1
	maxTitleLength = 50
)

type BookTitle struct {
	string
}

func NewBookTitle(title string) (BookTitle, error) {
	b := BookTitle{title}
	if err := b.validate(); err != nil {
		return BookTitle{}, err
	}

	return b, nil
}

func (b *BookTitle) validate() error {
	length := utf8.RuneCountInString(b.string)
	if length < minTitleLength || maxTitleLength < length {
		return newInvalidValueError(fmt.Errorf("invalid book title: got = %s", b.string))
	}

	return nil
}

func (b *BookTitle) String() string {
	return b.string
}
