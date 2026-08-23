package infra

import (
	"context"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/model"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/repository"
)

type bookQueryRepository struct{}

func NewBookQueryRepository() repository.BookQueryRepository {
	return &bookQueryRepository{}
}

type bookRow struct {
	BookID          string    `bun:"book_id"`
	Title           string    `bun:"title"`
	Summary         string    `bun:"summary"`
	Author          string    `bun:"author"`
	Price           uint      `bun:"price"`
	PublicationDate date.Date `bun:"publication_date"`
}

func (b *bookQueryRepository) List(ctx context.Context, executor database.Executor) (model.Books, error) {
	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to get bun connection: %w", err))
	}

	var rows []bookRow
	if err := bunDB.NewSelect().
		Column("book_id", "title", "summary", "author", "price", "publication_date").
		Table("book").
		Order("publication_date DESC").
		Scan(ctx, &rows); err != nil {
		return nil, repository.NewGetError(fmt.Errorf("failed to get book: %w", err))
	}

	books := make(model.Books, 0, len(rows))
	for _, v := range rows {
		books = append(books, toBookQueryModelFromBookRow(v))
	}

	return books, nil
}

func toBookQueryModelFromBookRow(row bookRow) model.Book {
	return model.Book{
		BookID:          row.BookID,
		Title:           row.Title,
		Summary:         row.Summary,
		Author:          row.Author,
		Price:           row.Price,
		PublicationDate: row.PublicationDate,
	}
}
