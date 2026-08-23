package openapitype

import (
	"fmt"

	"github.com/Sunaga0709/go-api-template/internal/date"

	oapitype "github.com/oapi-codegen/runtime/types"
)

func DateToOpenAPI(value date.Date) oapitype.Date {
	return oapitype.Date{Time: value.Time()}
}

func DateFromOpenAPI(value oapitype.Date) (date.Date, error) {
	d, err := date.Parse(value.String())
	if err != nil {
		return date.Date{}, newError(fmt.Errorf("failed to convert date: %w", err))
	}

	return d, nil
}
