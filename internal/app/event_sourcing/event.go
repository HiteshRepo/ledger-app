package event_sourcing

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/current_state"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
)

type Event interface {
	Apply(currState *current_state.CurrentState) (error, *order.Order, *order.Order)
	Display()
}
