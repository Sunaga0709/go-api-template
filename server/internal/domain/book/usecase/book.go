package usecase

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/book/repository"
)

type BookUsecase interface {
	Get(ctx context.Context, bookID string) (*model.Book, error)
	Create(ctx context.Context, input CreateInput) error
}

type CreateInput struct {
	Title           string
	Summary         string
	Author          string
	Price           uint
	PublicationDate date.Date
}

type bookUsecase struct {
	bookRepo         repository.BookRepository
	executorProvider database.ExecutorProvider
}

func NewBookUsecase(bookRepo repository.BookRepository, executorProvider database.ExecutorProvider) BookUsecase {
	return &bookUsecase{
		bookRepo:         bookRepo,
		executorProvider: executorProvider,
	}
}

func (b *bookUsecase) Get(ctx context.Context, bookID string) (*model.Book, error) {
	bid, err := model.NewBookID(bookID)
	if err != nil {
		return nil, newErrorFromModelError(err)
	}

	book, err := b.bookRepo.Get(ctx, b.executorProvider.Default(), bid)
	if err != nil {
		return nil, newErrorFromRepositoryError(err)
	}

	return book, nil
}

func (b *bookUsecase) Create(ctx context.Context, input CreateInput) error {
	bookID, err := model.GenerateBookID()
	if err != nil {
		return newErrorFromModelError(err)
	}
	title, err := model.NewBookTitle(input.Title)
	if err != nil {
		return newErrorFromModelError(err)
	}
	summary, err := model.NewBookSummary(input.Summary)
	if err != nil {
		return newErrorFromModelError(err)
	}
	author, err := model.NewBookAuthor(input.Author)
	if err != nil {
		return newErrorFromModelError(err)
	}
	price := model.NewBookPrice(input.Price)
	book := model.NewBook(
		bookID,
		title,
		summary,
		author,
		price,
		model.NewBookPublicationDate(input.PublicationDate),
	)

	if err := b.bookRepo.Create(ctx, b.executorProvider.Default(), book); err != nil {
		return newErrorFromRepositoryError(err)
	}

	return nil
}
