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

func TestLedgerRepository_Scenario1(t *testing.T) {
	id := uuid.New().String()
	name := "tomato"

	newProduct := product.NewProduct(id, name)

	err, matchDemand, matchSupply := newProduct.SupplyProduct(24, 100)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.SupplyProduct(20, 90)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.DemandProduct(22, 110)
	require.NoError(t, err)
	assert.NotNil(t, matchDemand)
	assert.NotNil(t, matchSupply)
	require.Equal(t, len(matchSupply), len(matchDemand))

	for i:=0; i<len(matchSupply); i++ {
		err = newProduct.TradeProduct(matchSupply[i], matchDemand[i])
		require.NoError(t, err)
	}

	err, matchDemand, matchSupply = newProduct.DemandProduct(21, 10)
	require.NoError(t, err)
	require.Nil(t, matchDemand)
	require.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.DemandProduct(21, 40)
	require.NoError(t, err)
	require.Nil(t, matchDemand)
	require.Nil(t, matchSupply)

	err, matchDemand, matchSupply = newProduct.SupplyProduct(19, 50)
	require.NoError(t, err)
	require.NotNil(t, matchDemand)
	require.NotNil(t, matchSupply)
	require.Equal(t, len(matchSupply), len(matchDemand))

	for i:=0; i<len(matchSupply); i++ {
		err = newProduct.TradeProduct(matchSupply[i], matchDemand[i])
		require.NoError(t, err)
	}

	repo := repository.NewWarehouseRepository()
	repo.Save(newProduct)

	expectedEvents := []event_sourcing.Event{
		event_sourcing.NewProductSupplyEvent(name, 24, 100),
		event_sourcing.NewProductSupplyEvent(name, 20, 90),
		event_sourcing.NewProductDemandEvent(name, 22, 110),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(20), Qty: decimal.NewFromFloat(90)},
			&order.Order{Price: decimal.NewFromFloat(20), Qty: decimal.NewFromFloat(90)}),
		event_sourcing.NewProductDemandEvent(name, 21, 10),
		event_sourcing.NewProductDemandEvent(name, 21, 40),
		event_sourcing.NewProductSupplyEvent(name, 19, 50),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(19), Qty: decimal.NewFromFloat(20)},
			&order.Order{Price: decimal.NewFromFloat(19), Qty: decimal.NewFromFloat(20)}),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(19), Qty: decimal.NewFromFloat(10)},
			&order.Order{Price: decimal.NewFromFloat(19), Qty: decimal.NewFromFloat(10)}),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(19), Qty: decimal.NewFromFloat(20)},
			&order.Order{Price: decimal.NewFromFloat(19), Qty: decimal.NewFromFloat(20)}),
	}

	actualProduct := repo.Get(id, name)
	actualEvents := actualProduct.GetEvents()

	require.Equal(t, len(expectedEvents), len(actualEvents))

	for i := 0; i < len(expectedEvents); i++ {
		assert.Equal(t, typeofobject(expectedEvents[i]), typeofobject(actualEvents[i]))
	}
}

func TestLedgerRepository_Scenario2(t *testing.T) {
	potatoId := uuid.New().String()
	potatoName := "potato"
	potato := product.NewProduct(potatoId, potatoName)
	tomatoId := uuid.New().String()
	tomatoName := "tomato"
	tomato := product.NewProduct(tomatoId, tomatoName)

	err, matchDemand, matchSupply := tomato.DemandProduct(110, 1)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = potato.DemandProduct(110, 10)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = tomato.DemandProduct(110, 10)
	require.NoError(t, err)
	assert.Nil(t, matchDemand)
	assert.Nil(t, matchSupply)

	err, matchDemand, matchSupply = potato.SupplyProduct(110, 1)
	require.NoError(t, err)
	require.NotNil(t, matchDemand)
	require.NotNil(t, matchSupply)
	require.Equal(t, len(matchSupply), len(matchDemand))

	for i:=0; i<len(matchSupply); i++ {
		err = potato.TradeProduct(matchSupply[i], matchDemand[i])
		require.NoError(t, err)
	}

	err, matchDemand, matchSupply = potato.SupplyProduct(110, 7)
	require.NoError(t, err)
	require.NotNil(t, matchDemand)
	require.NotNil(t, matchSupply)
	require.Equal(t, len(matchSupply), len(matchDemand))

	for i:=0; i<len(matchSupply); i++ {
		err = potato.TradeProduct(matchSupply[i], matchDemand[i])
		require.NoError(t, err)
	}

	err, matchDemand, matchSupply = potato.SupplyProduct(110, 2)
	require.NoError(t, err)
	require.NotNil(t, matchDemand)
	require.NotNil(t, matchSupply)
	require.Equal(t, len(matchSupply), len(matchDemand))

	for i:=0; i<len(matchSupply); i++ {
		err = potato.TradeProduct(matchSupply[i], matchDemand[i])
		require.NoError(t, err)
	}

	err, matchDemand, matchSupply = tomato.SupplyProduct(110, 11)
	require.NoError(t, err)
	require.NotNil(t, matchDemand)
	require.NotNil(t, matchSupply)
	require.Equal(t, len(matchSupply), len(matchDemand))

	for i:=0; i<len(matchSupply); i++ {
		err = tomato.TradeProduct(matchSupply[i], matchDemand[i])
		require.NoError(t, err)
	}

	repo := repository.NewWarehouseRepository()
	repo.Save(potato)
	repo.Save(tomato)

	expectedEventsPotato := []event_sourcing.Event{
		event_sourcing.NewProductDemandEvent(potatoName, 110, 10),
		event_sourcing.NewProductSupplyEvent(potatoName, 110, 1),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(1)},
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(1)}),
		event_sourcing.NewProductSupplyEvent(potatoName, 110, 7),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(7)},
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(7)}),
		event_sourcing.NewProductSupplyEvent(potatoName, 110, 2),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(2)},
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(2)}),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(1)},
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(1)}),
	}

	expectedEventsTomato := []event_sourcing.Event{
		event_sourcing.NewProductDemandEvent(tomatoName, 110, 1),
		event_sourcing.NewProductDemandEvent(tomatoName, 110, 10),
		event_sourcing.NewProductSupplyEvent(tomatoName, 110, 11),
		event_sourcing.NewTradeEvent(
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(11)},
			&order.Order{Price: decimal.NewFromFloat(110), Qty: decimal.NewFromFloat(11)}),
	}

	actualProductPotato := repo.Get(potatoId, potatoName)
	actualEventsPotato := actualProductPotato.GetEvents()

	actualProductTomato := repo.Get(tomatoId, tomatoName)
	actualEventsTomato := actualProductTomato.GetEvents()

	require.Equal(t, len(expectedEventsPotato), len(actualEventsPotato))
	for i := 0; i < len(expectedEventsPotato); i++ {
		assert.Equal(t, typeofobject(expectedEventsPotato[i]), typeofobject(actualEventsPotato[i]))
	}

	require.Equal(t, len(expectedEventsTomato), len(actualEventsTomato))
	for i := 0; i < len(expectedEventsTomato); i++ {
		assert.Equal(t, typeofobject(expectedEventsTomato[i]), typeofobject(actualEventsTomato[i]))
	}
}

func typeofobject(x interface{}) string {
	return fmt.Sprintf("%T", x)
}
