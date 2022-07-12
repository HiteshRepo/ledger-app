package event_sourcing

import "github.com/google/uuid"

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
