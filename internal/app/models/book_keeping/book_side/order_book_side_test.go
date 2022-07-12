package book_side_test

import (
	"github.com/google/uuid"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/book_side"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type orderBookSideSuite struct {
	suite.Suite
	bookSide book_side.BookSide
}

func TestBookSideSuite(t *testing.T) {
	suite.Run(t, new(orderBookSideSuite))
}

func (suite *orderBookSideSuite) SetupTest() {
	suite.bookSide = book_side.ProvideOrderBookSide()
	suite.bookSide.SetComparator(comparator.ProvideDemandComparator())
}

func (suite *orderBookSideSuite) TestGetOrdersOnEmptyOrderBook() {
	suite.Require().Empty(suite.bookSide.GetOrders())
}

func (suite *orderBookSideSuite) TestGetOrdersReturnsOneOrderAfterUpdate() {
	timeNow := time.Now().UnixNano()
	expected := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow},
	}
	err := suite.bookSide.UpdateOrders(expected)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, suite.bookSide.GetOrders())
}

func (suite *orderBookSideSuite) TestGetOrdersSortsAfterUpdate() {
	timeNow := time.Now().UnixNano()
	expected := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow},
	}
	err := suite.bookSide.UpdateOrders(expected)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, suite.bookSide.GetOrders())
}

func (suite *orderBookSideSuite) TestUpdatesTheOrderBookSideWithASingleIncrementalUpdateWithTheSamePriceAndNewQuantityForBid() {
	timeNow := time.Now().UnixNano()
	initialSnapshot := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	incrementalUpdate := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(5), Timestamp: timeNow},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(5), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	err := suite.bookSide.UpdateOrders(initialSnapshot)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, initialSnapshot, suite.bookSide.GetOrders())

	err = suite.bookSide.UpdateOrders(incrementalUpdate)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, expected, suite.bookSide.GetOrders())
}

func (suite *orderBookSideSuite) TestUpdatesTheOrderBookSideWithASingleIncrementalUpdateWithANewTimestamp() {
	initialTime := time.Now().Add(-1 * time.Hour).UnixNano()
	updatedTime := time.Now().UnixNano()

	initialSnapshot := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: initialTime},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: initialTime},
	}

	incrementalUpdate := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(5), Timestamp: updatedTime},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(5), Timestamp: updatedTime},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: initialTime},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: initialTime},
	}

	err := suite.bookSide.UpdateOrders(initialSnapshot)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, initialSnapshot, suite.bookSide.GetOrders())

	err = suite.bookSide.UpdateOrders(incrementalUpdate)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, expected, suite.bookSide.GetOrders())
}

func (suite *orderBookSideSuite) TestUpdatesTheOrderBookSideWithMultipleIncrementalUpdatesWithTheSamePriceAndADifferentQuantityForDemand() {
	timeNow := time.Now().UnixNano()
	initialSnapshot := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(8), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	incrementalUpdate := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(5), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(306), Qty: decimal.NewFromFloat(5), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(7), Timestamp: timeNow},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(306), Qty: decimal.NewFromFloat(5), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(8), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(7), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(5), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	err := suite.bookSide.UpdateOrders(initialSnapshot)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, initialSnapshot, suite.bookSide.GetOrders())

	err = suite.bookSide.UpdateOrders(incrementalUpdate)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, expected, suite.bookSide.GetOrders())
}

func (suite *orderBookSideSuite) TestRemovesAQuoteIfTheIncrementalUpdateHasAQuantityOfZeroForExistingId() {
	timeNow := time.Now().UnixNano()
	id := uuid.New().String()
	initialSnapshot := []*order.Order{
		{Id: id, Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(1), Timestamp: timeNow},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	incrementalUpdate := []*order.Order{
		{Id: id, Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(0), Timestamp: timeNow},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(3), Timestamp: timeNow},
	}

	err := suite.bookSide.UpdateOrders(initialSnapshot)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, initialSnapshot, suite.bookSide.GetOrders())

	err = suite.bookSide.UpdateOrders(incrementalUpdate)
	suite.Require().NoError(err)
	AssertEqualOrders(&suite.Suite, expected, suite.bookSide.GetOrders())
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
	suite.Assert().Equal(expected.Timestamp, actual.Timestamp)
}