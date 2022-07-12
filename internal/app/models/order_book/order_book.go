package order_book

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/book_side"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
)

type OrderBook interface {
	Get() (demands []*order.Order, supplies []*order.Order)
	Update(demands []*order.Order, supplies []*order.Order) error
}

type orderBook struct {
	demands  book_side.BookSide
	supplies book_side.BookSide
}

func ProvideOrderBook(demandComparator, supplyComparator comparator.Comparator) OrderBook {
	demandBookSide := book_side.ProvideOrderBookSide()
	demandBookSide.SetComparator(demandComparator)

	supplyBookSide := book_side.ProvideOrderBookSide()
	supplyBookSide.SetComparator(supplyComparator)

	return &orderBook{supplies: supplyBookSide, demands: demandBookSide}
}

func (m *orderBook) Get() (demands []*order.Order, supplies []*order.Order) {
	demands = m.demands.GetOrders()
	supplies = m.supplies.GetOrders()
	return
}

func (m *orderBook) Update(incomingDemands []*order.Order, incomingSupplies []*order.Order) error {
	err := m.demands.UpdateOrders(incomingDemands)
	if err != nil {
		return err
	}

	err = m.supplies.UpdateOrders(incomingSupplies)
	if err != nil {
		return err
	}

	return nil
}