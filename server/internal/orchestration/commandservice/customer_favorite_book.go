package commandservice

import (
	"context"

	bookuc "github.com/Sunaga0709/go-api-template/internal/domain/book/usecase"
	cfbuc "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/usecase"
	customeruc "github.com/Sunaga0709/go-api-template/internal/domain/customer/usecase"
)

type (
	CustomerFavoriteBookCommandService interface {
		Favorite(ctx context.Context, input FavoriteInput) error
	}

	FavoriteInput struct {
		CustomerID string
		BookID     string
	}
)

type customerFavoriteBookCommandService struct {
	customerUC             customeruc.CustomerUsecase
	bookUC                 bookuc.BookUsecase
	customerFavoriteBookUC cfbuc.CustomerFavoriteBookUsecase
}

func NewCustomerFavoriteBookCommandService(
	customerUC customeruc.CustomerUsecase,
	bookUC bookuc.BookUsecase,
	customerFavoriteBookUC cfbuc.CustomerFavoriteBookUsecase,
) CustomerFavoriteBookCommandService {
	return &customerFavoriteBookCommandService{
		customerUC:             customerUC,
		bookUC:                 bookUC,
		customerFavoriteBookUC: customerFavoriteBookUC,
	}
}

func (c *customerFavoriteBookCommandService) Favorite(ctx context.Context, input FavoriteInput) error {
	if _, err := c.customerUC.Get(ctx, input.CustomerID); err != nil {
		return newErrorFromCustomerUsecaseError(err)
	}

	if _, err := c.bookUC.Get(ctx, input.BookID); err != nil {
		return newErrorFromBookUsecaseError(err)
	}

	if err := c.customerFavoriteBookUC.Favorite(ctx, cfbuc.FavoriteInput{
		CustomerID: input.CustomerID,
		BookID:     input.BookID,
	}); err != nil {
		return newErrorFromCustomerFavoriteBookUsecaseError(err)
	}

	return nil
}
