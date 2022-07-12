package event_sourcing_test

import (
	"github.com/google/uuid"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/constants"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/event_sourcing"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/current_state"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order_book"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type productEventsSuite struct {
	suite.Suite
	currentState     *current_state.CurrentState
	existingDemands  []*order.Order
	existingSupplies []*order.Order
}

func TestProductEventsSuite(t *testing.T) {
	suite.Run(t, new(productEventsSuite))
}

func (suite *productEventsSuite) SetupTest() {
	suite.currentState = &current_state.CurrentState{
		OrderBook: order_book.ProvideOrderBook(comparator.ProvideDemandComparator(), comparator.ProvideSupplyComparator()),
	}

	timeNow := time.Now().UnixNano()

	suite.existingDemands = []*order.Order{
		{Id: uuid.New().String(), Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow},
		{Id: uuid.New().String(), Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow},
	}

	suite.existingSupplies = []*order.Order{
		{Id: uuid.New().String(), Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(7), Timestamp: timeNow},
		{Id: uuid.New().String(), Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	err := suite.currentState.OrderBook.Update(suite.existingDemands, suite.existingSupplies)
	suite.Require().NoError(err)
}

func (suite *productEventsSuite) TestProductSupplyEvent_ShouldAddToExistingSuppliesIfNoMatchTrade() {
	pse := event_sourcing.NewProductSupplyEvent("product-1", 500, 10)

	err, matchDemand, matchSupply := pse.Apply(suite.currentState)
	suite.Require().Error(err)
	suite.Assert().Contains(err.Error(), constants.OrderMismatchErrorMessage)
	suite.Assert().Nil(matchDemand)
	suite.Assert().Nil(matchSupply)

	expectedSupplies := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(7)},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(3)},
		{Price: decimal.NewFromFloat(500), Qty: decimal.NewFromFloat(10)},
	}

	expectedDemands := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11)},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10)},
	}

	existingDemands, newSupplies := suite.currentState.OrderBook.Get()

	AssertEqualOrders(&suite.Suite, existingDemands, expectedDemands)
	AssertEqualOrders(&suite.Suite, newSupplies, expectedSupplies)
}

func (suite *productEventsSuite) TestProductSupplyEvent_ShouldDecreaseFromMatchingDemandsIfMatchTrade() {
	pse := event_sourcing.NewProductSupplyEvent("product-1", 100, 10)

	err, matchDemand, matchSupply := pse.Apply(suite.currentState)
	suite.Require().NoError(err)
	suite.Assert().NotNil(matchDemand)
	suite.Assert().NotNil(matchSupply)

	expectedSupplies := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(7)},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(3)},
	}

	expectedDemands := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(1)},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10)},
	}

	newDemands, newSupplies := suite.currentState.OrderBook.Get()

	AssertEqualOrders(&suite.Suite, newDemands, expectedDemands)
	AssertEqualOrders(&suite.Suite, newSupplies, expectedSupplies)
}

func (suite *productEventsSuite) TestProductDemandEvent_ShouldAddToExistingDemandsIfNoMatchTrade() {
	pde := event_sourcing.NewProductDemandEvent("product-1", 99, 10)

	err, matchDemand, matchSupply := pde.Apply(suite.currentState)
	suite.Require().Error(err)
	suite.Assert().Contains(err.Error(), constants.OrderMismatchErrorMessage)
	suite.Assert().Nil(matchDemand)
	suite.Assert().Nil(matchSupply)

	expectedSupplies := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(7)},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(3)},
	}

	expectedDemands := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11)},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10)},
		{Price: decimal.NewFromFloat(99), Qty: decimal.NewFromFloat(10)},
	}

	existingDemands, newSupplies := suite.currentState.OrderBook.Get()

	AssertEqualOrders(&suite.Suite, existingDemands, expectedDemands)
	AssertEqualOrders(&suite.Suite, newSupplies, expectedSupplies)
}

func (suite *productEventsSuite) TestProductDemandEvent_ShouldRemoveFromExistingSuppliesIfMatchTrade() {
	pde := event_sourcing.NewProductDemandEvent("product-1", 100, 6)

	err, matchDemand, matchSupply := pde.Apply(suite.currentState)
	suite.Require().NoError(err)
	suite.Assert().NotNil(matchDemand)
	suite.Assert().NotNil(matchSupply)

	expectedSupplies := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(1)},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(3)},
	}

	expectedDemands := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11)},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10)},
	}

	newDemands, newSupplies := suite.currentState.OrderBook.Get()

	AssertEqualOrders(&suite.Suite, newDemands, expectedDemands)
	AssertEqualOrders(&suite.Suite, newSupplies, expectedSupplies)
}

func AssertEqualOrders(suite *suite.Suite, expected []*order.Order, actual []*order.Order) {
	suite.Assert().Equal(len(expected), len(actual))
	for i, q := range actual {
		AssertEqualOrder(suite, expected[i], q)
	}
}

func AssertEqualOrder(suite *suite.Suite, expected *order.Order, actual *order.Order) {
	suite.Assert().Equal(expected.Price, actual.Price)
	suite.Assert().Equal(expected.Qty, actual.Qty)
}
