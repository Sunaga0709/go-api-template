package repository

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/model"
)

type CustomerRepository interface {
	Get(ctx context.Context, executor database.Executor, customerID model.CustomerID) (*model.Customer, error)
	Create(ctx context.Context, executor database.Executor, customer *model.Customer) error
	Update(ctx context.Context, executor database.Executor, customer *model.Customer) error
}
