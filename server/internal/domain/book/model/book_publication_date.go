package model

import (
	"github.com/Sunaga0709/go-api-template/internal/date"
)

type BookPublicationDate struct {
	value date.Date
}

func NewBookPublicationDate(pd date.Date) BookPublicationDate {
	return BookPublicationDate{value: pd}
}

func (b *BookPublicationDate) Date() date.Date {
	return b.value
}
