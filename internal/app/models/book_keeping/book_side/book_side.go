package book_side

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
)

type BookSide interface {
	UpdateOrders(side []*order.Order) error
	GetOrders() []*order.Order
	SetComparator(comparator comparator.Comparator)
}
