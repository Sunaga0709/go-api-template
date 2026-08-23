package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/model"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/repository"
	"github.com/Sunaga0709/go-api-template/internal/order"
)

type customerQueryRepository struct{}

func NewCustomerQueryRepository() repository.CustomerQueryRepository {
	return &customerQueryRepository{}
}

func (c *customerQueryRepository) List(ctx context.Context, executor database.Executor, query repository.ListCustomerQuery) (model.Customers, error) {
	bun, err := database.UnwrapBun(executor)
	if err != nil {
		return nil, repository.NewUnknownError(err)
	}

	var rows []gen.Customer
	q := bun.NewSelect().
		Model(&rows)
	if query.Count > 0 {
		q = q.Limit(query.Count)
	}
	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}
	if query.Order == order.OrderAsc {
		q = q.Order("customer_id ASC")
	} else {
		q = q.Order("customer_id DESC")
	}

	if err := q.ExcludeColumn("created_at", "updated_at").Scan(ctx); err != nil {
		return nil, repository.NewGetError(fmt.Errorf("failed to get Customers: %w", err))
	}

	return convertCustomerSchemasToModels(rows)
}

func (c *customerQueryRepository) Get(ctx context.Context, executor database.Executor, query repository.GetCustomerQuery) (model.Customer, error) {
	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return model.Customer{}, repository.NewUnknownError(err)
	}

	var rows []customerDetailRow
	if scanErr := bunDB.NewSelect().
		Model(&rows).
		ModelTableExpr("customer AS c").
		ColumnExpr("c.customer_id AS customer_id").
		ColumnExpr("c.nickname AS nickname").
		ColumnExpr("c.birthday AS birthday").
		ColumnExpr("c.location AS location").
		ColumnExpr("b.book_id AS book_id").
		ColumnExpr("b.title AS title").
		ColumnExpr("b.summary AS summary").
		ColumnExpr("b.author AS author").
		ColumnExpr("b.price AS price").
		ColumnExpr("b.publication_date AS publication_date").
		Join(
			`LEFT JOIN (
				SELECT book_id
				FROM customer_favorite_book
				WHERE customer_id = ?
				ORDER BY created_at DESC
				LIMIT ?
				OFFSET ?
			) AS cfb ON TRUE`,
			query.CustomerID,
			query.Count,
			query.Offset,
		).
		Join("LEFT JOIN book AS b ON b.book_id = cfb.book_id").
		Where("c.customer_id = ?", query.CustomerID).
		Order("b.created_at DESC").
		Scan(ctx); scanErr != nil {
		return model.Customer{}, repository.NewGetError(fmt.Errorf("failed to get customer favorite bbok: %w", scanErr))
	}

	customer, convertErr := convertCustomerDetailRowsToModel(rows)
	if convertErr != nil {
		return model.Customer{}, repository.NewUnknownError(convertErr)
	}

	return customer, nil
}

type customerDetailRow struct {
	CustomerID      string     `bun:"customer_id"`
	Nickname        string     `bun:"nickname"`
	Birthday        date.Date  `bun:"birthday"`
	Location        string     `bun:"location"`
	BookID          *string    `bun:"book_id"`
	Title           *string    `bun:"title"`
	Summary         *string    `bun:"summary"`
	Author          *string    `bun:"author"`
	Price           *uint      `bun:"price"`
	PublicationDate *date.Date `bun:"publication_date"`
}

func convertCustomerSchemasToModels(rows []gen.Customer) (model.Customers, error) {
	customers := make(model.Customers, 0, len(rows))
	for _, v := range rows {
		year, month, day := v.Birthday.Date()
		birthday, err := date.NewDate(year, month, day)
		if err != nil {
			return nil, repository.NewUnknownError(fmt.Errorf("failed to convert birthday: %w", err))
		}

		customers = append(customers, &model.Customer{
			CustomerID: v.CustomerId,
			Nickname:   v.Nickname,
			Birthday:   birthday,
			Location:   v.Location,
		})
	}

	return customers, nil
}

func convertCustomerDetailRowsToModel(rows []customerDetailRow) (model.Customer, error) {
	if len(rows) == 0 {
		return model.Customer{}, repository.NewUnknownError(errors.New("failed to convert customer query model, because 0 row"))
	}

	var (
		customer      model.Customer
		isSetCustomer bool
	)
	books := make(model.Books, 0, len(rows))
	for _, v := range rows {
		if !isSetCustomer {
			isSetCustomer = true

			customer = model.Customer{
				CustomerID: v.CustomerID,
				Nickname:   v.Nickname,
				Birthday:   v.Birthday,
				Location:   v.Location,
			}
		}

		if v.BookID != nil && v.Title != nil && v.Summary != nil && v.Author != nil && v.Price != nil && v.PublicationDate != nil {
			books = append(books, model.Book{
				BookID:          *v.BookID,
				Title:           *v.Title,
				Summary:         *v.Summary,
				Author:          *v.Author,
				Price:           *v.Price,
				PublicationDate: *v.PublicationDate,
			})
		}
	}
	customer.FavoriteBooks = books

	return customer, nil
}
