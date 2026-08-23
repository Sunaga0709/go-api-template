package repository

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/model"
	"github.com/Sunaga0709/go-api-template/internal/order"
)

type CustomerQueryRepository interface {
	List(ctx context.Context, executor database.Executor, query ListCustomerQuery) (model.Customers, error)
	Get(ctx context.Context, executor database.Executor, query GetCustomerQuery) (model.Customer, error)
}

type ListCustomerQuery struct {
	Count  int
	Offset int
	Order  order.Order
}

type GetCustomerQuery struct {
	CustomerID string
	Count      int
	Offset     int
}
