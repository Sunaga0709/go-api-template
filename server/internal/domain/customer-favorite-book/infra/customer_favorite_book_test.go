package infra

import (
	"testing"

	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/model"
)

func TestCustomerFavoriteBookConversionsAndDiff(t *testing.T) {
	customerID, _ := model.NewCustomerID("customer")
	first, _ := model.NewBookID("first")
	second, _ := model.NewBookID("second")
	third, _ := model.NewBookID("third")
	row, err := toCustomerFavoriteBookSchemaFromModel(customerID, first)
	if err != nil || row.CustomerFavoriteBookId == "" || row.CustomerId != "customer" || row.BookId != "first" {
		t.Errorf("toCustomerFavoriteBookSchemaFromModel() = %#v, %v", row, err)
	}
	got, err := toCustomerFavoriteBookModelFromSchemas([]gen.CustomerFavoriteBook{{CustomerId: "customer", BookId: "first"}, {CustomerId: "customer", BookId: "second"}})
	if err != nil || got.CustomerID != customerID || len(got.BookIDs) != 2 {
		t.Errorf("toCustomerFavoriteBookModelFromSchemas() = %#v, %v", got, err)
	}
	if _, err := toCustomerFavoriteBookModelFromSchemas([]gen.CustomerFavoriteBook{{CustomerId: "", BookId: "first"}}); err == nil {
		t.Error("invalid customer error = nil")
	}
	if _, err := toCustomerFavoriteBookModelFromSchemas([]gen.CustomerFavoriteBook{{CustomerId: "customer", BookId: ""}}); err == nil {
		t.Error("invalid book error = nil")
	}
	expect := model.NewCustomerFavoriteBook(customerID, model.NewBookIDs([]model.BookID{first, third}))
	stored := model.NewCustomerFavoriteBook(customerID, model.NewBookIDs([]model.BookID{first, second}))
	diff := diffBookIDs(expect, stored)
	if len(diff.insert) != 1 || diff.insert[0] != third || len(diff.delete) != 1 || diff.delete[0] != second {
		t.Errorf("diffBookIDs() = %#v", diff)
	}
	if NewCustomerFavoriteBookRepository() == nil {
		t.Error("NewCustomerFavoriteBookRepository() = nil")
	}
}
