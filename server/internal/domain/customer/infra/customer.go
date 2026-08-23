package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/repository"
)

type customerRepository struct{}

func NewCustomerRepository() repository.CustomerRepository {
	return &customerRepository{}
}

func (c *customerRepository) Get(ctx context.Context, executor database.Executor, customerID model.CustomerID) (*model.Customer, error) {
	bun, err := database.UnwrapBun(executor)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to get bun connection: %w", err))
	}

	var row gen.Customer
	if err = bun.NewSelect().
		Model(&row).
		Where(fmt.Sprintf("%s = ?", gen.CustomerColumnCustomerId), customerID.String()).
		ExcludeColumn(gen.CustomerColumnCreatedAt, gen.CustomerColumnUpdatedAt).
		Limit(1).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewNotFoundError(fmt.Errorf("customer not found: %w", err))
		}

		return nil, repository.NewGetError(fmt.Errorf("failed to get customer: %w", err))
	}

	customer, err := toCustomerModelFromSchema(row)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *customerRepository) Create(ctx context.Context, executor database.Executor, customer *model.Customer) error {
	bun, err := database.UnwrapBun(executor)
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to get bun connection: %w", err))
	}

	row := toCustomerSchemaFromModel(customer)
	if _, err = bun.NewInsert().
		Model(&row).
		Exec(ctx); err != nil {
		return repository.NewCreateError(fmt.Errorf("failed to create customer: %w", err))
	}

	return nil
}

func (c *customerRepository) Update(ctx context.Context, executor database.Executor, customer *model.Customer) error {
	bun, err := database.UnwrapBun(executor)
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to get bun connection: %w", err))
	}

	row := toCustomerSchemaFromModel(customer)
	result, err := bun.NewUpdate().
		Model(&row).
		ExcludeColumn(gen.CustomerColumnCustomerId, gen.CustomerColumnCreatedAt).
		Where(fmt.Sprintf("%s = ?", gen.CustomerColumnCustomerId), row.CustomerId).
		Exec(ctx)
	if err != nil {
		return repository.NewUpdateError(fmt.Errorf("failed to update customer: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to get rows affected: %w", err))
	}
	if rowsAffected == 0 {
		return repository.NewNotFoundError(fmt.Errorf("customer not updated: customer id = %s", customer.CustomerID.String()))
	}

	return nil
}

func toCustomerSchemaFromModel(customer *model.Customer) gen.Customer {
	return gen.Customer{
		CustomerId:   customer.CustomerID.String(),
		Nickname:     customer.Nickname.String(),
		Birthday:     customer.Birthday.ToStdTime(),
		Location:     customer.Location.String(),
		RegisteredAt: customer.RegisteredAt.Time(),
	}
}

func toCustomerModelFromSchema(customer gen.Customer) (*model.Customer, error) {
	customerID, err := model.NewCustomerID(customer.CustomerId)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert customer id: %w", err))
	}
	nickname, err := model.NewCustomerNickname(customer.Nickname)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert customer nickname: %w", err))
	}
	birthday, err := model.NewCustomerBirthdayFromStdTime(customer.Birthday)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert customer birthday: %w", err))
	}
	location, err := model.NewCustomerLocationFromString(customer.Location)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert customer location: %w", err))
	}
	registeredAt := model.NewCustomerRegisteredAt(customer.RegisteredAt)

	c, err := model.NewCustomer(
		customerID,
		nickname,
		birthday,
		location,
		registeredAt,
	)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to convert customer: %w", err))
	}

	return c, nil
}
