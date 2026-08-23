package usecase

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/repository"
)

type (
	CustomerFavoriteBookUsecase interface {
		Favorite(ctx context.Context, input FavoriteInput) error
		Unfavorite(ctx context.Context, input UnfavoriteInput) error
	}

	FavoriteInput struct {
		CustomerID string
		BookID     string
	}

	UnfavoriteInput struct {
		CustomerID string
		BookID     string
	}
)

type customerFavoriteBookUsecase struct {
	customerFavoriteBookRepo repository.CustomerFavoriteBookRepository
	txManager                database.TxManager
}

func NewCustomerFavoriteBookUsecase(
	customerFavoriteBookRepo repository.CustomerFavoriteBookRepository,
	txManager database.TxManager,
) CustomerFavoriteBookUsecase {
	return &customerFavoriteBookUsecase{
		customerFavoriteBookRepo: customerFavoriteBookRepo,
		txManager:                txManager,
	}
}

func (c *customerFavoriteBookUsecase) Favorite(ctx context.Context, input FavoriteInput) error {
	customerID, err := model.NewCustomerID(input.CustomerID)
	if err != nil {
		return newErrorFromModelError(err)
	}
	bookID, err := model.NewBookID(input.BookID)
	if err != nil {
		return newErrorFromModelError(err)
	}

	if err := c.txManager.Run(ctx, func(ctx context.Context, executor database.TxExecutor) error {
		cfb, err := c.customerFavoriteBookRepo.GetByCustomerID(ctx, executor, customerID)
		if err != nil {
			return err
		}

		if err := cfb.Add(bookID); err != nil {
			return err
		}
		if err := c.customerFavoriteBookRepo.Update(ctx, executor, cfb); err != nil {
			return err
		}

		return nil
	}); err != nil {
		e := newError(err)
		if e.Kind == ErrorKindAlreadyFavorited {
			return e.WithClientMessage("この書籍は既にお気に入り済みです。")
		}

		return e
	}

	return nil
}

func (c *customerFavoriteBookUsecase) Unfavorite(ctx context.Context, input UnfavoriteInput) error {
	customerID, err := model.NewCustomerID(input.CustomerID)
	if err != nil {
		return newErrorFromModelError(err)
	}
	bookID, err := model.NewBookID(input.BookID)
	if err != nil {
		return newErrorFromModelError(err)
	}

	if err := c.txManager.Run(ctx, func(ctx context.Context, executor database.TxExecutor) error {
		cfb, err := c.customerFavoriteBookRepo.GetByCustomerID(ctx, executor, customerID)
		if err != nil {
			return err
		}

		if err := cfb.Remove(bookID); err != nil {
			return err
		}
		if err := c.customerFavoriteBookRepo.Update(ctx, executor, cfb); err != nil {
			return err
		}

		return nil
	}); err != nil {
		e := newError(err)
		if e.Kind == ErrorKindNotFavorited {
			return e.WithClientMessage("この書籍はお気に入りに登録されていません。")
		}

		return e
	}

	return nil
}
