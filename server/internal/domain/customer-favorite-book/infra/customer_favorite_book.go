package infra

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/repository"
)

type customerFavoriteBookRepository struct{}

func NewCustomerFavoriteBookRepository() repository.CustomerFavoriteBookRepository {
	return &customerFavoriteBookRepository{}
}

func (c *customerFavoriteBookRepository) GetByCustomerID(
	ctx context.Context,
	executor database.Executor,
	customerID model.CustomerID,
) (*model.CustomerFavoriteBook, error) {
	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return nil, repository.NewUnknownError(fmt.Errorf("failed to unwrap bun database connection: %w", err))
	}

	var rows []gen.CustomerFavoriteBook
	if err := bunDB.NewSelect().
		Model(&rows).
		Where(fmt.Sprintf("%s = ?", gen.CustomerFavoriteBookColumnCustomerId), customerID.String()).
		ExcludeColumn("customer_favorite_book_id", "created_at", "updated_at").
		ExcludeColumn(gen.CustomerFavoriteBookColumnCustomerFavoriteBookId, gen.CustomerFavoriteBookColumnCreatedAt, gen.CustomerFavoriteBookColumnUpdatedAt).
		Order(fmt.Sprintf("%s ASC", gen.CustomerFavoriteBookColumnCreatedAt)).
		Scan(ctx); err != nil {
		return nil, repository.NewGetError(fmt.Errorf("failed to get customer favorite book: %w", err))
	}

	return toCustomerFavoriteBookModelFromSchemas(rows)
}

func (c *customerFavoriteBookRepository) Update(ctx context.Context, executor database.Executor, cfb *model.CustomerFavoriteBook) error {
	customerID := cfb.CustomerID

	stored, err := c.GetByCustomerID(ctx, executor, customerID)
	if err != nil {
		return repository.NewUpdateError(fmt.Errorf("failed to get customer favorite book for update: %w", err))
	}

	diff := diffBookIDs(cfb, stored)

	if err := c.bulkCreateForUpdate(ctx, executor, customerID, diff.insert); err != nil {
		return err
	}
	if err := c.bulkDeleteForUpdate(ctx, executor, customerID, diff.delete); err != nil {
		return err
	}

	return nil
}

func (c *customerFavoriteBookRepository) bulkCreateForUpdate(ctx context.Context, executor database.Executor, customerID model.CustomerID, bookIDs []model.BookID) error {
	targetCount := len(bookIDs)
	if targetCount == 0 {
		return nil
	}

	cfbs := make([]gen.CustomerFavoriteBook, 0, targetCount)
	for _, v := range bookIDs {
		cfb, err := toCustomerFavoriteBookSchemaFromModel(customerID, v)
		if err != nil {
			return repository.NewUnknownError(fmt.Errorf("failed to convert customer favorite book schema for update: %w", err))
		}

		cfbs = append(cfbs, cfb)
	}

	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to unwrap bun database connection: %w", err))
	}

	result, err := bunDB.NewInsert().
		Model(&cfbs).
		Exec(ctx)
	if err != nil {
		return repository.NewUpdateError(fmt.Errorf("failed to create customer favorite book for update: %w", err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to get rows affected: %w", err))
	}
	if rowsAffected != int64(targetCount) {
		return repository.NewUpdateError(fmt.Errorf("unmatch created rows count for update: affected = %d, expect = %d", rowsAffected, targetCount))
	}

	return nil
}

func (c *customerFavoriteBookRepository) bulkDeleteForUpdate(ctx context.Context, executor database.Executor, customerID model.CustomerID, bookIDs []model.BookID) error {
	targetCount := len(bookIDs)
	if targetCount == 0 {
		return nil
	}

	bids := make([]string, 0, targetCount)
	for _, v := range bookIDs {
		bids = append(bids, v.String())
	}

	bunDB, err := database.UnwrapBun(executor)
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to unwrap bun database connection: %w", err))
	}

	result, err := bunDB.NewDelete().
		Model((*gen.CustomerFavoriteBook)(nil)).
		Where(fmt.Sprintf("%s = ?", gen.CustomerFavoriteBookColumnCustomerId), customerID.String()).
		Where(fmt.Sprintf("%s IN (?)", gen.CustomerFavoriteBookColumnBookId), bun.List(bids)).
		Exec(ctx)
	if err != nil {
		return repository.NewUpdateError(fmt.Errorf("failed to delete customer favorite book for update: %w", err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return repository.NewUnknownError(fmt.Errorf("failed to get rows affected: %w", err))
	}
	if rowsAffected != int64(targetCount) {
		return repository.NewUpdateError(fmt.Errorf("unmatch deleted rows count for update: affected = %d, expect = %d", rowsAffected, targetCount))
	}

	return nil
}

func toCustomerFavoriteBookSchemaFromModel(customerID model.CustomerID, bookID model.BookID) (gen.CustomerFavoriteBook, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return gen.CustomerFavoriteBook{}, repository.NewUnknownError(fmt.Errorf("failed to generate customer favorite book id: %w", err))
	}

	return gen.CustomerFavoriteBook{
		CustomerFavoriteBookId: id.String(),
		CustomerId:             customerID.String(),
		BookId:                 bookID.String(),
	}, nil
}

func toCustomerFavoriteBookModelFromSchemas(cfbs []gen.CustomerFavoriteBook) (*model.CustomerFavoriteBook, error) {
	var (
		isSetCustomerID bool
		customerID      model.CustomerID
	)
	bookIDs := make([]model.BookID, 0, len(cfbs))

	for _, v := range cfbs {
		if !isSetCustomerID {
			var err error
			customerID, err = model.NewCustomerID(v.CustomerId)
			if err != nil {
				return nil, repository.NewUnknownError(fmt.Errorf("failed to convert customer id: %w", err))
			}

			isSetCustomerID = true
		}

		bookID, err := model.NewBookID(v.BookId)
		if err != nil {
			return nil, repository.NewUnknownError(fmt.Errorf("failed to convert book id: %w", err))
		}

		bookIDs = append(bookIDs, bookID)
	}

	return model.NewCustomerFavoriteBook(customerID, model.NewBookIDs(bookIDs)), nil
}

type diffBookIDsResult struct {
	insert []model.BookID
	delete []model.BookID
}

func diffBookIDs(expect, stored *model.CustomerFavoriteBook) diffBookIDsResult {
	var insertBookIDs, deleteBookIDs []model.BookID

	// 追加対象
	for bookID := range expect.BookIDs {
		if _, exists := stored.BookIDs[bookID]; !exists {
			insertBookIDs = append(insertBookIDs, bookID)
		}
	}

	// 削除対象
	for bookID := range stored.BookIDs {
		if _, exists := expect.BookIDs[bookID]; !exists {
			deleteBookIDs = append(deleteBookIDs, bookID)
		}
	}

	return diffBookIDsResult{
		insert: insertBookIDs,
		delete: deleteBookIDs,
	}
}
