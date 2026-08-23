package model

import (
	"errors"
	"fmt"
)

type BookID struct {
	string
}

func NewBookID(bookID string) (BookID, error) {
	id := BookID{bookID}
	if err := id.validate(); err != nil {
		return BookID{}, err
	}

	return id, nil
}

func (b *BookID) validate() error {
	if b.string == "" {
		return newInvalidValueError(errors.New("book id is empty"))
	}

	return nil
}

func (b *BookID) String() string {
	return b.string
}

type BookIDs map[BookID]struct{}

func NewBookIDs(bookIDs []BookID) BookIDs {
	ids := make(map[BookID]struct{}, len(bookIDs))
	for _, v := range bookIDs {
		ids[v] = struct{}{}
	}

	return ids
}

func (b *BookIDs) Values() []BookID {
	ids := make([]BookID, 0, len(*b))
	for id := range *b {
		ids = append(ids, id)
	}

	return ids
}

func (b *BookIDs) add(bookID BookID) error {
	if _, exists := (*b)[bookID]; exists {
		return newAlreadyFavoritedError(fmt.Errorf("already favorite book: book id = %s", bookID.String()))
	}
	(*b)[bookID] = struct{}{}

	return nil
}

func (b *BookIDs) remove(bookID BookID) error {
	if _, exists := (*b)[bookID]; !exists {
		return newNotFavoritedError(fmt.Errorf("not favorited book: book id = %s", bookID.String()))
	}
	delete(*b, bookID)

	return nil
}
