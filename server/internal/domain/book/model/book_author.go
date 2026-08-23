package model

import (
	"fmt"
	"unicode/utf8"
)

const (
	minAuthorLength = 1
	maxAuthorLength = 30
)

type BookAuthor struct {
	string
}

func NewBookAuthor(author string) (BookAuthor, error) {
	b := BookAuthor{author}
	if err := b.validate(); err != nil {
		return BookAuthor{}, err
	}

	return b, nil
}

func (b *BookAuthor) validate() error {
	length := utf8.RuneCountInString(b.string)
	if length < minAuthorLength || maxAuthorLength < length {
		return newInvalidValueError(fmt.Errorf("invalid book author: got = %s", b.string))
	}

	return nil
}

func (b *BookAuthor) String() string {
	return b.string
}
