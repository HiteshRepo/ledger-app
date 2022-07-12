package book_keeping

import "github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"

type BookSide interface {
	UpdateOrders(side []*order.Order) error
	GetOrders() []*order.Order
	SetComparator(comparator Comparator)
}