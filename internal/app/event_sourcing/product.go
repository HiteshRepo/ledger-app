package event_sourcing

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/current_state"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order_book"
	"github.com/shopspring/decimal"
	"log"
	"strings"
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
	isOrderMatched, d, s := matchOrder(state.OrderBook, newSupplyOrder)

	if !isOrderMatched {
		return errors.New("order did not match"), nil, nil
	}

	return nil, d, s
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
	isOrderMatched, d, s := matchOrder(state.OrderBook, newDemandOrder)

	if !isOrderMatched {
		return errors.New("order did not match"), nil, nil
	}

	return nil, d, s
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

func matchOrder(orderbook order_book.OrderBook, o *order.Order) (bool, *order.Order, *order.Order) {
	zero := decimal.NewFromInt(0)
	demands, supplies := orderbook.Get()

	if strings.ToUpper(strings.TrimSpace(o.OrderType)) == "SUPPLY" {
		var maxDemand *order.Order
		for _, d := range demands {
			if d.Price.GreaterThanOrEqual(o.Price) && d.Qty.Equal(o.Qty) {
				if maxDemand == nil {
					maxDemand = d
					continue
				}

				if maxDemand.Price.LessThan(d.Price) {
					maxDemand = d
				}
			}
		}

		if maxDemand != nil {
			s := *o
			s.Qty = zero

			d := *maxDemand
			d.Qty = zero

			_ = orderbook.Update([]*order.Order{&d}, []*order.Order{&s})
			return true, &d, &s
		}
	}

	if strings.ToUpper(strings.TrimSpace(o.OrderType)) == "DEMAND" {
		var minSupply *order.Order
		for _, s := range supplies {
			if s.Price.LessThanOrEqual(o.Price) && s.Qty.Equal(o.Qty) {
				if minSupply == nil {
					minSupply = s
					continue
				}

				if minSupply.Price.GreaterThan(s.Price) {
					minSupply = s
				}
			}
		}

		if minSupply != nil {
			s := *minSupply
			s.Qty = zero

			d := *o
			d.Qty = zero

			_ = orderbook.Update([]*order.Order{&d}, []*order.Order{&s})
			return true, &d, &s
		}
	}

	return false, nil, nil
}