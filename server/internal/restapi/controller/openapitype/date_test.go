package openapitype

import (
	"testing"
	"time"

	oapitype "github.com/oapi-codegen/runtime/types"

	"github.com/Sunaga0709/go-api-template/internal/date"
)

func TestDateConversion(t *testing.T) {
	d, err := date.NewDate(2026, time.June, 27)
	if err != nil {
		t.Fatal(err)
	}
	openAPI := DateToOpenAPI(d)
	if openAPI.String() != "2026-06-27" {
		t.Errorf("DateToOpenAPI() = %s", openAPI.String())
	}
	got, err := DateFromOpenAPI(openAPI)
	if err != nil || got != d {
		t.Errorf("DateFromOpenAPI() = %v, %v", got, err)
	}
	if _, err := DateFromOpenAPI(oapitype.Date{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}); err != nil {
		t.Errorf("DateFromOpenAPI(min) error = %v", err)
	}
}
