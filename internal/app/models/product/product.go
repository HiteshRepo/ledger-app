package product

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/event_sourcing"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/current_state"
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
