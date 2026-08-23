package enumconvertor

import (
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func ConvertOrderToPrimitive(ord gen.Order) string {
	return string(ord)
}

func ConvertOrderToOpenAPI(ord string) gen.Order {
	switch ord {
	case string(gen.OrderDesc):
		return gen.OrderDesc
	default:
		return gen.OrderAsc
	}
}
