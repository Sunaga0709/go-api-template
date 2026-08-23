package queryservice

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/model"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/repository"
	"github.com/Sunaga0709/go-api-template/internal/order"
)

const (
	defaultListCustomerCount  = 100
	defaultListCustomerOffset = 0

	defaultGetCustomerFavoriteProductCount  = 100
	defaultGetCustomerFavoriteProductOffset = 0
)

type (
	CustomerQueryService interface {
		List(ctx context.Context, input ListCustomerInput) (ListCustomerOutput, error)
		Get(ctx context.Context, input GetCustomerInput) (GetCustomerOutput, error)
	}

	ListCustomerInput struct {
		Count  int
		Offset int
		Order  string
	}

	ListCustomerOutput struct {
		Customers model.Customers
	}

	GetCustomerInput struct {
		CustomerID string
		Count      *int
		Offset     *int
	}

	GetCustomerOutput struct {
		CustomerID    string
		Nickname      string
		Birthday      string
		Location      string
		FavoriteBooks []GetCustomerOutputBook
	}

	GetCustomerOutputBook struct {
		BookID          string
		Title           string
		Summary         string
		Author          string
		Price           uint
		PublicationDate date.Date
	}
)

type customerQueryService struct {
	customerQueryRepo repository.CustomerQueryRepository
	executorProvider  database.ExecutorProvider
}

func NewCustomerQueryService(
	customerQueryRepo repository.CustomerQueryRepository,
	executorProvider database.ExecutorProvider,
) CustomerQueryService {
	return &customerQueryService{
		customerQueryRepo: customerQueryRepo,
		executorProvider:  executorProvider,
	}
}

func (c *customerQueryService) List(ctx context.Context, input ListCustomerInput) (ListCustomerOutput, error) {
	count, offset, ord := defaultListCustomerCount, defaultListCustomerOffset, order.OrderAsc
	if input.Count > 0 {
		count = input.Count
	}
	if input.Offset > 0 {
		offset = input.Offset
	}
	// 正常にパースできた場合のみ、入力値を使用する
	if o, err := order.ParseOrder(input.Order); err == nil {
		ord = o
	}

	customers, err := c.customerQueryRepo.List(ctx, c.executorProvider.Default(), repository.ListCustomerQuery{
		Count:  count,
		Offset: offset,
		Order:  ord,
	})
	if err != nil {
		return ListCustomerOutput{}, newErrorFromQueryRepositoryError(err)
	}

	return ListCustomerOutput{Customers: customers}, nil
}

func (c *customerQueryService) Get(ctx context.Context, input GetCustomerInput) (GetCustomerOutput, error) {
	count := defaultGetCustomerFavoriteProductCount
	offset := defaultGetCustomerFavoriteProductOffset
	if input.Count != nil {
		count = *input.Count
	}
	if input.Offset != nil {
		offset = *input.Offset
	}

	customer, err := c.customerQueryRepo.Get(ctx, c.executorProvider.Default(), repository.GetCustomerQuery{
		CustomerID: input.CustomerID,
		Count:      count,
		Offset:     offset,
	})
	if err != nil {
		return GetCustomerOutput{}, newErrorFromQueryRepositoryError(err)
	}

	books := make([]GetCustomerOutputBook, 0, len(customer.FavoriteBooks))
	for _, v := range customer.FavoriteBooks {
		books = append(books, GetCustomerOutputBook{
			BookID:          v.BookID,
			Title:           v.Title,
			Summary:         v.Summary,
			Author:          v.Author,
			Price:           v.Price,
			PublicationDate: v.PublicationDate,
		})
	}

	return GetCustomerOutput{
		CustomerID:    customer.CustomerID,
		Nickname:      customer.Nickname,
		Birthday:      customer.Birthday.String(),
		Location:      customer.Location,
		FavoriteBooks: books,
	}, nil
}
