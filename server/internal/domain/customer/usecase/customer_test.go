package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/model"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/repository"
)

type testProvider struct{ executor database.Executor }

func (p testProvider) Default() database.Executor { return p.executor }

type testCustomerRepository struct {
	customer                     *model.Customer
	getErr, createErr, updateErr error
	created, updated             *model.Customer
}

func (r *testCustomerRepository) Get(context.Context, database.Executor, model.CustomerID) (*model.Customer, error) {
	return r.customer, r.getErr
}
func (r *testCustomerRepository) Create(_ context.Context, _ database.Executor, c *model.Customer) error {
	r.created = c
	return r.createErr
}
func (r *testCustomerRepository) Update(_ context.Context, _ database.Executor, c *model.Customer) error {
	r.updated = c
	return r.updateErr
}

func TestCustomerUsecase(t *testing.T) {
	id := uuid.Must(uuid.NewV7()).String()
	cid, _ := model.NewCustomerID(id)
	nickname, _ := model.NewCustomerNickname("old")
	birthdayDate, _ := date.NewDate(2000, time.January, 1)
	location, _ := model.NewCustomerLocationFromString("UTC")
	customer, _ := model.NewCustomer(cid, nickname, model.NewCustomerBirthday(birthdayDate), location, model.DefaultCustomerRegisteredAt())
	repo := &testCustomerRepository{customer: customer}
	uc := NewCustomerUsecase(repo, testProvider{database.NewBunExecutor(nil)})
	if got, err := uc.Get(context.Background(), id); err != nil || got != customer {
		t.Fatalf("Get() = %v, %v", got, err)
	}
	if _, err := uc.Get(context.Background(), "bad"); err == nil {
		t.Error("Get(invalid) error = nil")
	}
	repo.getErr = repository.NewGetError(errors.New("db"))
	if _, err := uc.Get(context.Background(), id); err == nil {
		t.Error("Get(repo error) error = nil")
	}
	repo.getErr = nil
	output, err := uc.Create(context.Background())
	if err != nil || repo.created == nil || output.CustomerID.String() == "" {
		t.Fatalf("Create() = %#v, %v", output, err)
	}
	repo.createErr = repository.NewCreateError(errors.New("db"))
	if _, err := uc.Create(context.Background()); err == nil {
		t.Error("Create(repo error) error = nil")
	}
	repo.createErr = nil
	input := UpdateInput{CustomerID: id, Nickname: "new", Birthday: "2001-02-03", Location: "Asia/Tokyo"}
	if err := uc.Update(context.Background(), input); err != nil || repo.updated == nil || repo.updated.Nickname.String() != "new" {
		t.Fatalf("Update() error = %v, updated = %#v", err, repo.updated)
	}
	for _, invalid := range []UpdateInput{{CustomerID: "bad"}, {CustomerID: id, Nickname: ""}, {CustomerID: id, Nickname: "ok", Birthday: "bad"}, {CustomerID: id, Nickname: "ok", Birthday: "2001-02-03", Location: "bad"}} {
		if err := uc.Update(context.Background(), invalid); err == nil {
			t.Errorf("Update(%#v) error = nil", invalid)
		}
	}
	repo.getErr = repository.NewGetError(errors.New("db"))
	if err := uc.Update(context.Background(), input); err == nil {
		t.Error("Update(get error) error = nil")
	}
	repo.getErr = nil
	repo.updateErr = repository.NewUpdateError(errors.New("db"))
	if err := uc.Update(context.Background(), input); err == nil {
		t.Error("Update(repo error) error = nil")
	}
}
