package repository

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/model"
)

type CustomerFavoriteBookRepository interface {
	GetByCustomerID(ctx context.Context, executor database.Executor, customerID model.CustomerID) (*model.CustomerFavoriteBook, error)
	Update(ctx context.Context, executor database.Executor, cfb *model.CustomerFavoriteBook) error
}
