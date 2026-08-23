package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/repository"
)

type testTxManager struct{ runErr error }

func (m testTxManager) Run(ctx context.Context, fn func(context.Context, database.TxExecutor) error) error {
	if m.runErr != nil {
		return m.runErr
	}
	return fn(ctx, nil)
}
func (m testTxManager) WithTxExecutor(database.TxExecutor) database.TxManager { return m }

type testFavoriteRepository struct {
	favorite          *model.CustomerFavoriteBook
	getErr, updateErr error
	updated           bool
}

func (r *testFavoriteRepository) GetByCustomerID(context.Context, database.Executor, model.CustomerID) (*model.CustomerFavoriteBook, error) {
	return r.favorite, r.getErr
}
func (r *testFavoriteRepository) Update(context.Context, database.Executor, *model.CustomerFavoriteBook) error {
	r.updated = true
	return r.updateErr
}

func TestCustomerFavoriteBookUsecase(t *testing.T) {
	cid, _ := model.NewCustomerID("customer")
	repo := &testFavoriteRepository{favorite: model.NewCustomerFavoriteBook(cid, model.NewBookIDs(nil))}
	uc := NewCustomerFavoriteBookUsecase(repo, testTxManager{})
	if err := uc.Favorite(context.Background(), FavoriteInput{CustomerID: "customer", BookID: "book"}); err != nil || !repo.updated {
		t.Fatalf("Favorite() error = %v", err)
	}
	if err := uc.Favorite(context.Background(), FavoriteInput{CustomerID: "", BookID: "book"}); err == nil {
		t.Error("Favorite(invalid) error = nil")
	}
	if err := uc.Favorite(context.Background(), FavoriteInput{CustomerID: "customer", BookID: "book"}); err == nil {
		t.Error("Favorite(existing) error = nil")
	}
	if err := uc.Unfavorite(context.Background(), UnfavoriteInput{CustomerID: "customer", BookID: "book"}); err != nil {
		t.Fatalf("Unfavorite() error = %v", err)
	}
	if err := uc.Unfavorite(context.Background(), UnfavoriteInput{CustomerID: "customer", BookID: "book"}); err == nil {
		t.Error("Unfavorite(missing) error = nil")
	}
	repo.getErr = repository.NewGetError(errors.New("db"))
	if err := uc.Favorite(context.Background(), FavoriteInput{CustomerID: "customer", BookID: "other"}); err == nil {
		t.Error("Favorite(get error) error = nil")
	}
	repo.getErr = nil
	repo.updateErr = repository.NewUpdateError(errors.New("db"))
	if err := uc.Favorite(context.Background(), FavoriteInput{CustomerID: "customer", BookID: "other"}); err == nil {
		t.Error("Favorite(update error) error = nil")
	}
}
