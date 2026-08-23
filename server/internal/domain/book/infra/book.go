package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/repository"
)

type bookRepository struct{}

func NewBookRepository() repository.BookRepository {
	return &bookRepository{}
}

func (b *bookRepository) Get(ctx context.Context, executor database.Executor, bookID model.BookID) (*model.Book, error) {
	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to unwrap bun database connection: %w", err))
	}

	var row gen.Book
	if err := bunDB.NewSelect().
		Model(&row).
		ExcludeColumn(gen.BookColumnCreatedAt, gen.BookColumnUpdatedAt).
		Where(fmt.Sprintf("%s = ?", gen.BookColumnBookId), bookID.String()).
		Limit(1).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewNotFoundError(fmt.Errorf("book not found (book id = %s)", bookID.String()))
		}

		return nil, repository.NewGetError(fmt.Errorf("failed to get book: %w", err))
	}

	return toBookModelFromSchema(row)
}

func (b *bookRepository) Create(ctx context.Context, executor database.Executor, book *model.Book) error {
	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to unwrap bun database connection: %w", err))
	}

	row := toBookSchemaFromModel(book)
	if _, err := bunDB.NewInsert().Model(&row).Exec(ctx); err != nil {
		return repository.NewCreateError(fmt.Errorf("failed to create Book: %w", err))
	}

	return nil
}

func toBookSchemaFromModel(book *model.Book) gen.Book {
	return gen.Book{
		BookId:          book.BookID.String(),
		Title:           book.Title.String(),
		Summary:         book.Summary.String(),
		Author:          book.Author.String(),
		Price:           book.Price.Uint(),
		PublicationDate: time.Date(book.PublicationDate.Date().Year(), book.PublicationDate.Date().Month(), book.PublicationDate.Date().Day(), 0, 0, 0, 0, time.UTC),
	}
}

func toBookModelFromSchema(book gen.Book) (*model.Book, error) {
	bookID, err := model.NewBookID(book.BookId)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert book id: %w", err))
	}
	title, err := model.NewBookTitle(book.Title)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert book title: %w", err))
	}
	summary, err := model.NewBookSummary(book.Summary)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert book summary: %w", err))
	}
	author, err := model.NewBookAuthor(book.Author)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert book author: %w", err))
	}
	price := model.NewBookPrice(book.Price)
	publicationDateValue, err := date.NewDate(book.PublicationDate.Year(), book.PublicationDate.Month(), book.PublicationDate.Day())
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert book publication date: %w", err))
	}
	publicationDate := model.NewBookPublicationDate(publicationDateValue)

	return model.NewBook(
		bookID,
		title,
		summary,
		author,
		price,
		publicationDate,
	), nil
}
