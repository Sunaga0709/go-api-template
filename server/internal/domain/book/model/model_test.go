package model

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Sunaga0709/go-api-template/internal/date"
)

func TestNewBookID(t *testing.T) {
	valid := uuid.Must(uuid.NewV7()).String()
	tests := []struct {
		name, id string
		wantErr  bool
	}{
		{name: "valid UUID", id: valid},
		{name: "empty value", wantErr: true},
		{name: "malformed value", id: "not-a-uuid", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewBookID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.String() != tt.id {
				t.Errorf("NewBookID() = %q, want %q", got.String(), tt.id)
			}
		})
	}
}

func TestGenerateBookID(t *testing.T) {
	got, err := GenerateBookID()
	if err != nil {
		t.Fatalf("GenerateBookID() error = %v", err)
	}
	if _, err := uuid.Parse(got.String()); err != nil {
		t.Errorf("GenerateBookID() = %q, not UUID: %v", got.String(), err)
	}
}

func TestBookTextValues(t *testing.T) {
	tests := []struct {
		name           string
		new            func(string) (string, error)
		valid, tooLong string
	}{
		{"title", func(v string) (string, error) { x, e := NewBookTitle(v); return x.String(), e }, "a", strings.Repeat("あ", maxTitleLength+1)},
		{"summary", func(v string) (string, error) { x, e := NewBookSummary(v); return x.String(), e }, "", strings.Repeat("あ", maxSummaryLength+1)},
		{"author", func(v string) (string, error) { x, e := NewBookAuthor(v); return x.String(), e }, "a", strings.Repeat("あ", maxAuthorLength+1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range []struct {
				name, value string
				wantErr     bool
			}{
				{"normal", tt.valid, false}, {"maximum rune length", strings.TrimSuffix(tt.tooLong, "あ"), false}, {"over maximum", tt.tooLong, true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					got, err := tt.new(tc.value)
					if (err != nil) != tc.wantErr {
						t.Fatalf("constructor error = %v, wantErr %v", err, tc.wantErr)
					}
					if !tc.wantErr && got != tc.value {
						t.Errorf("String() = %q, want %q", got, tc.value)
					}
				})
			}
			if tt.name != "summary" {
				if _, err := tt.new(""); err == nil {
					t.Error("empty value error = nil, want error")
				}
			}
		})
	}
}

func TestBookValueObjects(t *testing.T) {
	d, err := date.NewDate(2026, time.June, 27)
	if err != nil {
		t.Fatal(err)
	}
	price := NewBookPrice(99)
	if price.Uint() != 99 || price.Int() != 99 {
		t.Errorf("BookPrice accessors = %v, %v", price.Uint(), price.Int())
	}
	published := NewBookPublicationDate(d)
	if published.Date() != d {
		t.Errorf("Date() = %v, want %v", published.Date(), d)
	}
	id, _ := NewBookID(uuid.Must(uuid.NewV7()).String())
	title, _ := NewBookTitle("title")
	summary, _ := NewBookSummary("summary")
	author, _ := NewBookAuthor("author")
	book := NewBook(id, title, summary, author, price, published)
	if book.BookID != id || book.Title != title || book.Summary != summary || book.Author != author || book.Price != price || book.PublicationDate != published {
		t.Errorf("NewBook() = %#v, fields were not retained", book)
	}
}

func TestError(t *testing.T) {
	err := errors.New("cause")
	for _, tt := range []struct {
		kind ErrorKind
		want string
	}{{ErrorKindUnknown, "unknown"}, {ErrorKindInvalidValue, "invalid_value"}, {ErrorKind(99), "unknown"}} {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
	if got := newError(err, ErrorKindUnknown).Error(); !strings.Contains(got, "unknown") || !strings.Contains(got, "cause") {
		t.Errorf("Error() = %q", got)
	}
	if newInvalidValueError(err).Kind != ErrorKindInvalidValue {
		t.Error("newInvalidValueError kind mismatch")
	}
	if newUnknownError(err).Kind != ErrorKindUnknown {
		t.Error("newUnknownError kind mismatch")
	}
}
