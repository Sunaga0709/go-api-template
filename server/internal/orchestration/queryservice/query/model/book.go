package model

import "github.com/Sunaga0709/go-api-template/internal/date"

type Books []Book

type Book struct {
	BookID          string
	Title           string
	Summary         string
	Author          string
	Price           uint
	PublicationDate date.Date
}
