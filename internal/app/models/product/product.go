package product

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/constants"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/event_sourcing"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/current_state"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order_book"
)

type Product struct {
	Id           string
	name         string
	events       []event_sourcing.Event
	currentState *current_state.CurrentState
}

func NewProduct(id string, name string) *Product {
	return &Product{
		Id:           id,
		name:         name,
		currentState: &current_state.CurrentState{OrderBook: order_book.ProvideOrderBook(comparator.ProvideDemandComparator(), comparator.ProvideSupplyComparator())},
	}
}

func (p *Product) SupplyProduct(price, quantity float64) (error, *order.Order, *order.Order) {
	ev := event_sourcing.NewProductSupplyEvent(p.name, price, quantity)

	err, matchDemand, matchSupply := p.AddEvent(ev)
	if err != nil {
		return err, nil, nil
	}

	return nil, matchDemand, matchSupply
}

func (p *Product) DemandProduct(price, quantity float64) (error, *order.Order, *order.Order) {
	ev := event_sourcing.NewProductDemandEvent(p.name, price, quantity)

	err, matchDemand, matchSupply := p.AddEvent(ev)
	if err != nil {
		return err, nil, nil
	}

	return nil, matchDemand, matchSupply
}

func (p *Product) TradeProduct(matchSupply, matchDemand *order.Order) error {
	ev := event_sourcing.NewTradeEvent(matchSupply, matchDemand)
	err, _, _ := p.AddEvent(ev)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) GetCurrentState() *current_state.CurrentState {
	return p.currentState
}

func (p *Product) AddEvent(ev event_sourcing.Event) (error, *order.Order, *order.Order) {
	err, matchDemand, matchSupply := ev.Apply(p.currentState)

	if err != nil && err.Error() != constants.OrderMismatchErrorMessage {
		return err, nil, nil
	}

	p.events = append(p.events, ev)
	return nil, matchDemand, matchSupply
}

func (p *Product) GetEvents() []event_sourcing.Event {
	return p.events
}
