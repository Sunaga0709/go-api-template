package infra

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/model"
)

func TestBookSchemaConversions(t *testing.T) {
	id, _ := model.NewBookID(uuid.Must(uuid.NewV7()).String())
	title, _ := model.NewBookTitle("title")
	summary, _ := model.NewBookSummary("summary")
	author, _ := model.NewBookAuthor("author")
	d, _ := date.NewDate(2026, time.June, 27)
	book := model.NewBook(id, title, summary, author, model.NewBookPrice(99), model.NewBookPublicationDate(d))
	row := toBookSchemaFromModel(book)
	if row.BookId != id.String() || row.Price != 99 || !row.PublicationDate.Equal(time.Date(2026, time.June, 27, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("toBookSchemaFromModel() = %#v", row)
	}
	got, err := toBookModelFromSchema(row)
	if err != nil || got.BookID != id || got.Title != title || got.PublicationDate.Date() != d {
		t.Errorf("toBookModelFromSchema() = %#v, %v", got, err)
	}
	for _, row := range []gen.Book{{BookId: "bad"}, {BookId: id.String(), Title: ""}, {BookId: id.String(), Title: "ok", Summary: string(make([]byte, 501)), Author: "a"}, {BookId: id.String(), Title: "ok", Summary: "", Author: "", PublicationDate: time.Now()}} {
		if _, err := toBookModelFromSchema(row); err == nil {
			t.Errorf("toBookModelFromSchema(%#v) error = nil", row)
		}
	}
	if NewBookRepository() == nil {
		t.Error("NewBookRepository() = nil")
	}
}
