package comparator

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
)

func ProvideDemandComparator() Comparator {
	return func(order1 *order.Order, order2 *order.Order) (int, error) {
		if order1.Price.Equal((*order2).Price) && order1.Qty.Equal((*order2).Qty) {
			return cmp(order1.Timestamp, order2.Timestamp), nil
		}

		if order1.Price.Equal((*order2).Price) {
			return order2.Qty.Cmp((*order1).Qty), nil
		}

		return order1.Price.Cmp((*order2).Price), nil
	}
}
