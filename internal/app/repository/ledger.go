package repository

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/event_sourcing"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/product"
)

type LedgerRepository struct {
	inMemoryLedger map[string][]event_sourcing.Event
}

func NewWarehouseRepository() *LedgerRepository {
	return &LedgerRepository{inMemoryLedger: make(map[string][]event_sourcing.Event)}
}

func (wr *LedgerRepository) Get(id string, name string) *product.Product {
	newProduct := product.NewProduct(id, name)

	if events, ok := wr.inMemoryLedger[id]; ok {
		for _, e := range events {
			_, _, _ = newProduct.AddEvent(e)
		}
	}

	return newProduct
}

func (wr *LedgerRepository) Save(product *product.Product) {
	wr.inMemoryLedger[product.Id] = product.GetEvents()
}
