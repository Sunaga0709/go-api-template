package model

import (
	"fmt"
	"unicode/utf8"
)

const (
	minSummaryLength = 0
	maxSummaryLength = 500
)

type BookSummary struct {
	string
}

func NewBookSummary(summary string) (BookSummary, error) {
	b := BookSummary{summary}
	if err := b.validate(); err != nil {
		return BookSummary{}, err
	}

	return b, nil
}

func (b *BookSummary) validate() error {
	length := utf8.RuneCountInString(b.string)
	if length < minSummaryLength || maxSummaryLength < length {
		return newInvalidValueError(fmt.Errorf("invalid book summary: got = %s", b.string))
	}

	return nil
}

func (b *BookSummary) String() string {
	return b.string
}
