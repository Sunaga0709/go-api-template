package usecase

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/repository"
)

type (
	CustomerUsecase interface {
		Get(ctx context.Context, customerID string) (*model.Customer, error)
		Create(ctx context.Context) (CreateOutput, error)
		Update(ctx context.Context, input UpdateInput) error
	}

	CreateOutput struct {
		CustomerID model.CustomerID
	}

	UpdateInput struct {
		CustomerID string
		Nickname   string
		Birthday   string
		Location   string
	}
)

type customerUsecase struct {
	customerRepo     repository.CustomerRepository
	executorProvider database.ExecutorProvider
}

func NewCustomerUsecase(
	customerRepo repository.CustomerRepository,
	executorProvider database.ExecutorProvider,
) CustomerUsecase {
	return &customerUsecase{
		customerRepo:     customerRepo,
		executorProvider: executorProvider,
	}
}

func (c *customerUsecase) Get(ctx context.Context, customerID string) (*model.Customer, error) {
	cid, err := model.NewCustomerID(customerID)
	if err != nil {
		return nil, newErrorFromModelError(err)
	}

	customer, err := c.customerRepo.Get(ctx, c.executorProvider.Default(), cid)
	if err != nil {
		return nil, newErrorFromRepositoryError(err)
	}

	return customer, nil
}

func (c *customerUsecase) Create(ctx context.Context) (CreateOutput, error) {
	customerID, err := model.GenerateCustomerID()
	if err != nil {
		return CreateOutput{}, newErrorFromModelError(err)
	}

	customer, err := model.NewCustomer(customerID, model.DefaultCustomerNickname(), model.DefaultCustomerBirthday(), model.DefaultCustomerLocation(), model.DefaultCustomerRegisteredAt())
	if err != nil {
		return CreateOutput{}, newErrorFromModelError(err)
	}

	if err := c.customerRepo.Create(ctx, c.executorProvider.Default(), customer); err != nil {
		return CreateOutput{}, newErrorFromRepositoryError(err)
	}

	return CreateOutput{
		CustomerID: customer.CustomerID,
	}, nil
}

func (c *customerUsecase) Update(ctx context.Context, input UpdateInput) error {
	cid, err := model.NewCustomerID(input.CustomerID)
	if err != nil {
		return newErrorFromModelError(err)
	}
	nickname, err := model.NewCustomerNickname(input.Nickname)
	if err != nil {
		return newErrorFromModelError(err)
	}
	bd, err := date.Parse(input.Birthday)
	if err != nil {
		return newInvalidValueError(err)
	}
	birthday := model.NewCustomerBirthday(bd)
	location, err := model.NewCustomerLocationFromString(input.Location)
	if err != nil {
		return newErrorFromModelError(err)
	}

	customer, err := c.customerRepo.Get(ctx, c.executorProvider.Default(), cid)
	if err != nil {
		return newErrorFromRepositoryError(err)
	}

	if err := customer.Update(nickname, birthday, location); err != nil {
		return newErrorFromModelError(err)
	}
	if err := c.customerRepo.Update(ctx, c.executorProvider.Default(), customer); err != nil {
		return newErrorFromRepositoryError(err)
	}

	return nil
}
