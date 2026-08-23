package model

import (
	"errors"
	"strings"
	"testing"
)

func TestCustomerIDAndBookID(t *testing.T) {
	for _, tt := range []struct {
		name string
		new  func(string) (string, error)
	}{
		{"customer", func(v string) (string, error) { x, e := NewCustomerID(v); return x.String(), e }},
		{"book", func(v string) (string, error) { x, e := NewBookID(v); return x.String(), e }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.new("id"); err != nil || got != "id" {
				t.Errorf("constructor() = %q, %v", got, err)
			}
			if _, err := tt.new(""); err == nil {
				t.Error("constructor(empty) error = nil")
			}
		})
	}
}

func TestBookIDs(t *testing.T) {
	a, _ := NewBookID("a")
	b, _ := NewBookID("b")
	ids := NewBookIDs([]BookID{a, a})
	if len(ids) != 1 {
		t.Errorf("NewBookIDs() length = %d, want 1", len(ids))
	}
	if err := ids.add(b); err != nil {
		t.Fatalf("add() error = %v", err)
	}
	if err := ids.add(b); err == nil {
		t.Error("add(existing) error = nil")
	}
	if err := ids.remove(b); err != nil {
		t.Fatalf("remove() error = %v", err)
	}
	if err := ids.remove(b); err == nil {
		t.Error("remove(missing) error = nil")
	}
	values := ids.Values()
	if len(values) != 1 || values[0] != a {
		t.Errorf("Values() = %v", values)
	}
}

func TestCustomerFavoriteBook(t *testing.T) {
	customerID, _ := NewCustomerID("customer")
	bookID, _ := NewBookID("book")
	favorite := NewCustomerFavoriteBook(customerID, NewBookIDs(nil))
	if favorite.CustomerID != customerID {
		t.Error("NewCustomerFavoriteBook() customer id mismatch")
	}
	if err := favorite.Add(bookID); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if err := favorite.Remove(bookID); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
}

func TestError(t *testing.T) {
	err := errors.New("cause")
	for _, tt := range []struct {
		kind ErrorKind
		want string
	}{{ErrorKindUnknown, errorKindUnknownString}, {ErrorKindInvalidValue, errorKindInvalidValueString}, {ErrorKindAlreadyFavorited, errorKindAlreadyFavoritedString}, {ErrorKindNotFavorited, errorKindNotFavoritedString}, {ErrorKind(99), errorKindUnknownString}} {
		if got := tt.kind.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
	if newInvalidValueError(err).Kind != ErrorKindInvalidValue || newAlreadyFavoritedError(err).Kind != ErrorKindAlreadyFavorited || newNotFavoritedError(err).Kind != ErrorKindNotFavorited {
		t.Error("error constructor kind mismatch")
	}
	if got := newError(ErrorKindUnknown, err).Error(); !strings.Contains(got, "unknown") || !strings.Contains(got, "cause") {
		t.Errorf("Error() = %q", got)
	}
}
