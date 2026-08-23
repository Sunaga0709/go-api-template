package repository

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/model"
)

type BookQueryRepository interface {
	List(ctx context.Context, executor database.Executor) (model.Books, error)
}
