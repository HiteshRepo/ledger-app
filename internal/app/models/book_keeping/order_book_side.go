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

func (a *orderBookSide) GetOrders() []*order.Order {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	return a.orders
}

func (a *orderBookSide) SetComparator(comparator Comparator) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.comparator = comparator
}