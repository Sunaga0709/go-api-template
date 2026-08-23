package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/repository"
)

type testProvider struct{ executor database.Executor }

func (p testProvider) Default() database.Executor { return p.executor }

type testBookRepository struct {
	got               *model.Book
	getErr, createErr error
	created           *model.Book
}

func (r *testBookRepository) Get(context.Context, database.Executor, model.BookID) (*model.Book, error) {
	return r.got, r.getErr
}
func (r *testBookRepository) Create(_ context.Context, _ database.Executor, b *model.Book) error {
	r.created = b
	return r.createErr
}

func TestBookUsecase(t *testing.T) {
	id := uuid.Must(uuid.NewV7()).String()
	bookID, _ := model.NewBookID(id)
	title, _ := model.NewBookTitle("title")
	summary, _ := model.NewBookSummary("summary")
	author, _ := model.NewBookAuthor("author")
	book := model.NewBook(bookID, title, summary, author, model.NewBookPrice(1), model.NewBookPublicationDate(date.Min()))
	repo := &testBookRepository{got: book}
	uc := NewBookUsecase(repo, testProvider{database.NewBunExecutor(nil)})
	got, err := uc.Get(context.Background(), id)
	if err != nil || got != book {
		t.Fatalf("Get() = %v, %v", got, err)
	}
	if _, err := uc.Get(context.Background(), "invalid"); err == nil {
		t.Error("Get(invalid ID) error = nil")
	}
	repo.getErr = repository.NewGetError(errors.New("db"))
	if _, err := uc.Get(context.Background(), id); err == nil {
		t.Error("Get(repository error) error = nil")
	}
	repo.getErr = nil
	input := CreateInput{Title: "title", Summary: "summary", Author: "author", Price: 99, PublicationDate: date.Min()}
	if err := uc.Create(context.Background(), input); err != nil || repo.created == nil {
		t.Fatalf("Create() error = %v, created = %v", err, repo.created)
	}
	if err := uc.Create(context.Background(), CreateInput{Title: ""}); err == nil {
		t.Error("Create(invalid input) error = nil")
	}
	repo.createErr = repository.NewCreateError(errors.New("db"))
	if err := uc.Create(context.Background(), input); err == nil {
		t.Error("Create(repository error) error = nil")
	}
	_ = time.Second
}
