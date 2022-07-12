package repository_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/event_sourcing"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/product"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLedgerRepository_Get(t *testing.T) {
	id := uuid.New().String()
	name := "tomato"

	newProduct := product.NewProduct(id, name)

	err, matchDemand, matchSupply := newProduct.SupplyProduct(500, 11)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.SupplyProduct(100, 20)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.DemandProduct(200, 15)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.DemandProduct(100, 20)
	require.NoError(t, err)
	require.NotNil(t, matchDemand)
	require.NotNil(t, matchSupply)

	err = newProduct.TradeProduct(matchSupply, matchDemand)
	require.NoError(t, err)

	repo := repository.NewWarehouseRepository()
	repo.Save(newProduct)

	expectedEvents := []event_sourcing.Event{
		event_sourcing.NewProductSupplyEvent(name, 500, 11),
		event_sourcing.NewProductSupplyEvent(name, 100, 20),
		event_sourcing.NewProductDemandEvent(name, 200, 15),
		event_sourcing.NewProductDemandEvent(name, 100, 20),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(20)},
			&order.Order{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(20)}),
	}

	actualProduct := repo.Get(id, name)
	actualEvents := actualProduct.GetEvents()

	assert.Len(t, expectedEvents, len(actualEvents))

	for i:=0; i<len(expectedEvents); i++ {
		assert.Equal(t, typeofobject(expectedEvents[i]), typeofobject(actualEvents[i]))
	}
}

func typeofobject(x interface{}) string {
	return  fmt.Sprintf("%T", x)
}

