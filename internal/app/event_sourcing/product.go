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