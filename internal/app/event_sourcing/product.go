package event_sourcing

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/constants"
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

func (pse productSupplyEvent) Apply(state *current_state.CurrentState) (error, []*order.Order, []*order.Order) {
	newSupplyOrder := &order.Order{
		Id:        uuid.New().String(),
		Price:     decimal.NewFromFloat(pse.price),
		Qty:       decimal.NewFromFloat(pse.qty),
		OrderType: constants.SupplyOrderType,
	}

	_ = state.OrderBook.Update(nil, []*order.Order{newSupplyOrder})
	d, s := matchOrder(state.OrderBook, newSupplyOrder)

	if len(d) == 0 && len(s) == 0 {
		return errors.New(constants.OrderMismatchErrorMessage), nil, nil
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

func (pde productDemandEvent) Apply(state *current_state.CurrentState) (error, []*order.Order, []*order.Order) {
	newDemandOrder := &order.Order{
		Id:        uuid.New().String(),
		Price:     decimal.NewFromFloat(pde.price),
		Qty:       decimal.NewFromFloat(pde.qty),
		OrderType: constants.DemandOrderType,
	}

	_ = state.OrderBook.Update([]*order.Order{newDemandOrder}, nil)
	d, s := matchOrder(state.OrderBook, newDemandOrder)

	if len(d) == 0 && len(s) == 0 {
		return errors.New(constants.OrderMismatchErrorMessage), nil, nil
	}

	return nil, d, s
}

func (pde productDemandEvent) Display() {
	log.Printf("Demand order for product (%s) registered with quantity: %v, status: %s at %d\n", pde.productName, pde.qty, pde.status, pde.timestamp)
}

type tradeEvent struct {
	id     uuid.UUID
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

func (te tradeEvent) Apply(_ *current_state.CurrentState) (error, []*order.Order, []*order.Order) {
	return nil, nil, nil
}

func (te tradeEvent) Display() {
	log.Printf("Trade occured with supply id: %v and demand id: %v", te.supply.Id, te.demand.Id)
}

func matchOrder(orderbook order_book.OrderBook, o *order.Order) ([]*order.Order, []*order.Order) {
	zero := decimal.NewFromInt(0)
	currentOrder := *o

	matchDemands := make([]*order.Order, 0)
	matchSupplies := make([]*order.Order, 0)

	var matchSupply, matchDemand *order.Order

	if strings.ToUpper(strings.TrimSpace(currentOrder.OrderType)) == constants.SupplyOrderType {
		isFullFillPossible := true
		for isFullFillPossible && currentOrder.Qty.GreaterThan(zero) {
			var maxDemand *order.Order
			demands, _ := orderbook.Get()
			for _, d := range demands {
				if d.Price.GreaterThanOrEqual(currentOrder.Price) {
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
				isFullFillPossible, matchDemand, matchSupply = fulfillOrder(orderbook, maxDemand, &currentOrder, zero)
				if isFullFillPossible {
					matchDemands = append(matchDemands, matchDemand)
					matchSupplies = append(matchSupplies, matchSupply)
					currentOrder.Qty = currentOrder.Qty.Sub(matchDemand.Qty)
				}
			} else {
				isFullFillPossible = false
			}
		}
	}

	if strings.ToUpper(strings.TrimSpace(currentOrder.OrderType)) == constants.DemandOrderType {
		isFullFillPossible := true
		for isFullFillPossible && currentOrder.Qty.GreaterThan(zero) {
			_, supplies := orderbook.Get()
			var minSupply *order.Order
			for _, s := range supplies {
				if s.Price.LessThanOrEqual(currentOrder.Price) {
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
				isFullFillPossible, matchDemand, matchSupply = fulfillOrder(orderbook, &currentOrder, minSupply, zero)
				if isFullFillPossible {
					matchDemands = append(matchDemands, matchDemand)
					matchSupplies = append(matchSupplies, matchSupply)
					currentOrder.Qty = currentOrder.Qty.Sub(matchSupply.Qty)
				}
			} else {
				isFullFillPossible = false
			}
		}
	}

	return matchDemands, matchSupplies
}

func fulfillOrder(orderbook order_book.OrderBook, demand *order.Order, supply *order.Order, zero decimal.Decimal) (bool, *order.Order, *order.Order) {
	s := *supply
	d := *demand

	sq, _ := s.Qty.Float64()
	dq, _ := d.Qty.Float64()

	updatedSupplyQty := decimal.NewFromFloat(sq - dq)
	if updatedSupplyQty.IsNegative() || updatedSupplyQty.IsZero() {
		updatedSupplyQty = zero
	}
	newSupply := &order.Order{Id: s.Id, Price: s.Price, Qty: updatedSupplyQty, OrderType: constants.SupplyOrderType, Timestamp: s.Timestamp}

	updatedDemandQty := decimal.NewFromFloat(dq - sq)
	if updatedDemandQty.IsNegative() || updatedDemandQty.IsZero() {
		updatedDemandQty = zero
	}

	s.Qty = zero
	d.Qty = zero
	newDemand := &order.Order{Id: d.Id, Price: d.Price, Qty: updatedDemandQty, OrderType: constants.DemandOrderType, Timestamp: d.Timestamp}

	_ = orderbook.Update([]*order.Order{&d}, []*order.Order{&s})
	_ = orderbook.Update([]*order.Order{newDemand}, []*order.Order{newSupply})

	fullFilledQty := min(supply.Qty, demand.Qty)

	matchDemand := &order.Order{
		Id:        d.Id,
		Price:     d.Price,
		Qty:       fullFilledQty,
		OrderType: constants.DemandOrderType,
		Timestamp: d.Timestamp,
	}

	matchSupply := &order.Order{
		Id:        s.Id,
		Price:     s.Price,
		Qty:       fullFilledQty,
		OrderType: constants.SupplyOrderType,
		Timestamp: s.Timestamp,
	}

	return true, matchDemand, matchSupply
}


func max(v1, v2 decimal.Decimal) decimal.Decimal {
	if v1.GreaterThan(v2) {
		return v1
	}
	return v2
}

func min(v1, v2 decimal.Decimal) decimal.Decimal {
	if v1.LessThan(v2) {
		return v1
	}
	return v2
}