package queryservice

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/model"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/repository"
)

type (
	BookQueryService interface {
		List(ctx context.Context) (model.Books, error)
	}
)

type bookQueryService struct {
	bookQuery        repository.BookQueryRepository
	executorProvider database.ExecutorProvider
}

func NewBookQueryService(bookQuery repository.BookQueryRepository, executorProvider database.ExecutorProvider) BookQueryService {
	return &bookQueryService{
		bookQuery:        bookQuery,
		executorProvider: executorProvider,
	}
}

func (b *bookQueryService) List(ctx context.Context) (model.Books, error) {
	books, err := b.bookQuery.List(ctx, b.executorProvider.Default())
	if err != nil {
		return nil, newErrorFromQueryRepositoryError(err)
	}

	return books, nil
}
