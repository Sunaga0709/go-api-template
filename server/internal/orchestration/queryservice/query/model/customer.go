package model

import (
	"github.com/Sunaga0709/go-api-template/internal/date"
)

type Customers []*Customer

type Customer struct {
	CustomerID    string
	Nickname      string
	Birthday      date.Date
	Location      string
	FavoriteBooks Books
}
