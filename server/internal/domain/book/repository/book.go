package repository

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/model"
)

type BookRepository interface {
	Get(ctx context.Context, executor database.Executor, bookID model.BookID) (*model.Book, error)
	Create(ctx context.Context, executor database.Executor, book *model.Book) error
}
