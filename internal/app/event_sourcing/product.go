package event_sourcing

import (
	"github.com/google/uuid"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/current_state"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/shopspring/decimal"
	"log"
)

type productSupplyEvent struct {
	id          uuid.UUID
	productName string
	price       float64
	qty         float64
	status      string
	timestamp   int64
}

func NewProductSupplyEvent(productName string, price, quantity float64) Event {
	return productSupplyEvent{
		id:          uuid.New(),
		productName: productName,
		price:       price,
		qty:         quantity,
	}
}

func (pse productSupplyEvent) Apply(state *current_state.CurrentState) (error, *order.Order, *order.Order) {
	newSupplyOrder := &order.Order{
		Id:        uuid.New().String(),
		Price:     decimal.NewFromFloat(pse.price),
		Qty:       decimal.NewFromFloat(pse.qty),
		OrderType: "supply",
		Status:    "pending",
	}

	_ = state.OrderBook.Update(nil, []*order.Order{newSupplyOrder})

	// match order

	return nil, nil, nil
}

func (pse productSupplyEvent) Display() {
	log.Printf("Supply order for product (%s) registered with quantity: %v, status: %s at %d\n", pse.productName, pse.qty, pse.status, pse.timestamp)
}

type productDemandEvent struct {
	id          uuid.UUID
	productName string
	price       float64
	qty         float64
	status      string
	timestamp   int64
}

func NewProductDemandEvent(productName string, price, quantity float64) Event {
	return productDemandEvent{
		id:          uuid.New(),
		productName: productName,
		price:       price,
		qty:         quantity,
	}
}

func (pde productDemandEvent) Apply(state *current_state.CurrentState) (error, *order.Order, *order.Order) {
	newDemandOrder := &order.Order{
		Id:        uuid.New().String(),
		Price:     decimal.NewFromFloat(pde.price),
		Qty:       decimal.NewFromFloat(pde.qty),
		OrderType: "demand",
		Status:    "pending",
	}

	_ = state.OrderBook.Update([]*order.Order{newDemandOrder}, nil)

	// match order

	return nil, nil, nil
}

func (pde productDemandEvent) Display() {
	log.Printf("Demand order for product (%s) registered with quantity: %v, status: %s at %d\n", pde.productName, pde.qty, pde.status, pde.timestamp)
}

type tradeEvent struct {
	id          uuid.UUID
	supply *order.Order
	demand *order.Order
}

func NewTradeEvent(supplyEvent *order.Order, demandEvent *order.Order) Event {
	return tradeEvent{
		id:     uuid.New(),
		supply: supplyEvent,
		demand: demandEvent,
	}
}

func (te tradeEvent) Apply(_ *current_state.CurrentState) (error, *order.Order, *order.Order) {
	return nil, nil, nil
}

func (te tradeEvent) Display() {
	log.Printf("Trade occured with supply id: %v and demand id: %v", te.supply.Id, te.demand.Id)
}