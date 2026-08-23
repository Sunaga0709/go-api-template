package order

import "fmt"

type Order int

const (
	OrderAsc Order = iota
	OrderDesc

	orderAscString  = "asc"
	orderDescString = "desc"
)

func ParseOrder(order string) (Order, error) {
	switch order {
	case orderAscString:
		return OrderAsc, nil
	case orderDescString:
		return OrderDesc, nil
	default:
		return 0, newError(fmt.Errorf("invalid order string: got = %s", order))
	}
}
