package book_keeping

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
)

type Comparator func(order1 *order.Order, order2 *order.Order) (int, error)

func cmp(t1, t2 int64) int {
	if t1 > t2 { return -1 }
	if t1 < t2 { return 1 }
	return 0
}