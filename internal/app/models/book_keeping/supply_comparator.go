package book_keeping

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
)

func ProvideSupplyComparator() Comparator {
	return func(order1 *order.Order, order2 *order.Order) (int, error) {
		if order1.Price.Equal((*order2).Price) && order1.Qty.Equal((*order2).Qty) {
			return cmp(order1.Timestamp, order2.Timestamp), nil
		}

		if order1.Price.Equal((*order2).Price) {
			return order1.Qty.Cmp((*order2).Qty), nil
		}

		return order2.Price.Cmp((*order1).Price), nil
	}
}
