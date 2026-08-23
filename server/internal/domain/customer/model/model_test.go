package model

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Sunaga0709/go-api-template/internal/date"
	"github.com/Sunaga0709/go-api-template/internal/location"
)

func TestCustomerID(t *testing.T) {
	id := uuid.Must(uuid.NewV7()).String()
	got, err := NewCustomerID(id)
	if err != nil || got.String() != id {
		t.Fatalf("NewCustomerID() = %v, %v", got, err)
	}
	if _, invalidErr := NewCustomerID("invalid"); invalidErr == nil {
		t.Error("NewCustomerID(invalid) error = nil")
	}
	generated, err := GenerateCustomerID()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := uuid.Parse(generated.String()); err != nil {
		t.Errorf("GenerateCustomerID() = %q: %v", generated.String(), err)
	}
}

func TestCustomerNickname(t *testing.T) {
	for _, tt := range []struct {
		name, value string
		wantErr     bool
	}{{"minimum", "a", false}, {"maximum runes", strings.Repeat("あ", maxNicknameLength), false}, {"empty", "", true}, {"too long", strings.Repeat("あ", maxNicknameLength+1), true}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCustomerNickname(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewCustomerNickname() error = %v", err)
			}
			if !tt.wantErr && got.String() != tt.value {
				t.Errorf("String() = %q", got.String())
			}
		})
	}
	defaultNicknameValue := DefaultCustomerNickname()
	if got := defaultNicknameValue.String(); got != defaultNickname {
		t.Errorf("DefaultCustomerNickname() = %q", got)
	}
}

func TestCustomerBirthday(t *testing.T) {
	d, _ := date.NewDate(2000, time.February, 29)
	birthday := NewCustomerBirthday(d)
	if birthday.Date() != d || !birthday.IsPast() {
		t.Errorf("NewCustomerBirthday() = %#v", birthday)
	}
	if got := birthday.ToStdTime(); !got.Equal(time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ToStdTime() = %v", got)
	}
	fromTime, err := NewCustomerBirthdayFromStdTime(time.Date(2026, time.June, 27, 23, 0, 0, 0, time.FixedZone("x", 9*3600)))
	if err != nil || fromTime.Date().String() != "2026-06-27" {
		t.Errorf("NewCustomerBirthdayFromStdTime() = %#v, %v", fromTime, err)
	}
	defaultBirthday := DefaultCustomerBirthday()
	if !defaultBirthday.IsDefault() {
		t.Error("DefaultCustomerBirthday().IsDefault() = false")
	}
}

func TestCustomerLocation(t *testing.T) {
	for _, tt := range []struct {
		name, value string
		wantErr     bool
	}{{"UTC", "UTC", false}, {"Japan", "Asia/Tokyo", false}, {"invalid", "Not/AZone", true}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCustomerLocationFromString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewCustomerLocationFromString() error = %v", err)
			}
			if !tt.wantErr && got.String() != tt.value {
				t.Errorf("String() = %q", got.String())
			}
		})
	}
	if got := NewCustomerLocation(nil); got.String() != location.JST().String() || got.Location() != location.JST() {
		t.Errorf("NewCustomerLocation(nil) = %v", got)
	}
	var zero CustomerLocation
	if zero.String() != location.JST().String() || zero.Location() != location.JST() {
		t.Errorf("zero CustomerLocation = %v", zero)
	}
}

func TestCustomerRegisteredAt(t *testing.T) {
	when := time.Date(2026, time.June, 27, 12, 0, 0, 0, time.FixedZone("x", 9*3600))
	registered := NewCustomerRegisteredAt(when)
	if !registered.IsRegistered() || registered.Time().Location() != time.UTC {
		t.Errorf("NewCustomerRegisteredAt() = %#v", registered)
	}
	zero := DefaultCustomerRegisteredAt()
	if zero.IsRegistered() {
		t.Error("DefaultCustomerRegisteredAt() registered")
	}
	zero.Register()
	if !zero.IsRegistered() || zero.Time().Location() != time.UTC {
		t.Error("Register() did not set UTC timestamp")
	}
}

func TestCustomer(t *testing.T) {
	id, _ := NewCustomerID(uuid.Must(uuid.NewV7()).String())
	nickname, _ := NewCustomerNickname("name")
	past, _ := date.NewDate(2000, time.January, 1)
	birthday := NewCustomerBirthday(past)
	loc, _ := NewCustomerLocationFromString("UTC")
	customer, err := NewCustomer(id, nickname, birthday, loc, DefaultCustomerRegisteredAt())
	if err != nil || customer.IsRegistered() {
		t.Fatalf("NewCustomer() = %#v, %v", customer, err)
	}
	if updateErr := customer.Update(nickname, birthday, loc); updateErr != nil {
		t.Fatalf("Update() error = %v", updateErr)
	}
	if !customer.IsRegistered() || customer.Nickname != nickname || customer.Location != loc {
		t.Errorf("Update() = %#v", customer)
	}
	registeredDefaultBirthday, err := NewCustomer(id, nickname, DefaultCustomerBirthday(), loc, NewCustomerRegisteredAt(time.Now()))
	if err != nil || !registeredDefaultBirthday.IsRegistered() {
		t.Errorf("NewCustomer() registered default birthday = %#v, %v", registeredDefaultBirthday, err)
	}
}

func TestError(t *testing.T) {
	err := errors.New("cause")
	for _, tt := range []struct {
		kind ErrorKind
		want string
	}{{ErrorKindUnknown, errorKindUnknown}, {ErrorKindInvalidValue, errorKindInvalidValue}, {ErrorKind(99), errorKindUnknown}} {
		if tt.kind.String() != tt.want {
			t.Errorf("String() = %q, want %q", tt.kind.String(), tt.want)
		}
	}
	if newInvalidValueError(err).Kind != ErrorKindInvalidValue || newUnknownError(err).Kind != ErrorKindUnknown {
		t.Error("error constructor kind mismatch")
	}
	if got := newError(err, ErrorKindUnknown).Error(); !strings.Contains(got, "unknown") || !strings.Contains(got, "cause") {
		t.Errorf("Error() = %q", got)
	}
}
