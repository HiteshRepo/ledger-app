package book_keeping

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"sync"
)

func ProvideOrderBookSide() BookSide {
	return &orderBookSide{orders: make([]*order.Order, 0)}
}

type orderBookSide struct {
	mtx sync.RWMutex

	orders     []*order.Order
	comparator Comparator
}
