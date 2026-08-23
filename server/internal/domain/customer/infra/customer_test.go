package infra

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Sunaga0709/go-api-template/internal/database/gen"
	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/domain/customer/model"
)

func TestCustomerSchemaConversions(t *testing.T) {
	id, _ := model.NewCustomerID(uuid.Must(uuid.NewV7()).String())
	nickname, _ := model.NewCustomerNickname("name")
	d, _ := date.NewDate(2000, time.January, 2)
	loc, _ := model.NewCustomerLocationFromString("Asia/Tokyo")
	registered := model.NewCustomerRegisteredAt(time.Date(2026, time.June, 27, 0, 0, 0, 0, time.UTC))
	customer, _ := model.NewCustomer(id, nickname, model.NewCustomerBirthday(d), loc, registered)
	row := toCustomerSchemaFromModel(customer)
	if row.CustomerId != id.String() || row.Nickname != "name" || row.Location != "Asia/Tokyo" {
		t.Errorf("toCustomerSchemaFromModel() = %#v", row)
	}
	got, err := toCustomerModelFromSchema(row)
	if err != nil || got.CustomerID != id || got.Nickname != nickname || got.Location.String() != "Asia/Tokyo" {
		t.Errorf("toCustomerModelFromSchema() = %#v, %v", got, err)
	}
	for _, row := range []gen.Customer{{CustomerId: "bad"}, {CustomerId: id.String(), Nickname: "", Birthday: time.Now(), Location: "UTC"}, {CustomerId: id.String(), Nickname: "ok", Birthday: time.Now(), Location: "bad"}} {
		if _, err := toCustomerModelFromSchema(row); err == nil {
			t.Errorf("toCustomerModelFromSchema(%#v) error = nil", row)
		}
	}
	if NewCustomerRepository() == nil {
		t.Error("NewCustomerRepository() = nil")
	}
}
