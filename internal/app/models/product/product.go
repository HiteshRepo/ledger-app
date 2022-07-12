package product

import (
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


func (p *Product) GetCurrentState() *current_state.CurrentState {
	return p.currentState
}

func (p *Product) AddEvent(ev event_sourcing.Event) (error, *order.Order, *order.Order) {
	err, matchDemand, matchSupply := ev.Apply(p.currentState)

	if err != nil && err.Error() != "order did not match" {
		return err, nil, nil
	}

	p.events = append(p.events, ev)
	return nil, matchDemand, matchSupply
}

func (p *Product) GetEvents() []event_sourcing.Event {
	return p.events
}