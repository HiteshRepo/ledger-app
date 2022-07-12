package comparator_test

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"sort"
	"testing"
	"time"
)

type supplyComparatorSuite struct {
	suite.Suite
	comparator comparator.Comparator
}

func TestSupplyComparatorSuite(t *testing.T) {
	suite.Run(t, new(supplyComparatorSuite))
}

func (suite *supplyComparatorSuite) SetupTest() {
	suite.comparator = comparator.ProvideSupplyComparator()
}

func (suite *supplyComparatorSuite) TestSortsByPriceInAscending() {
	timeNow := time.Now()

	input := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow.UnixNano()},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
	}

	sort.Slice(input, func(i, j int) bool {
		comparison, err := suite.comparator(input[i], input[j])
		suite.Require().NoError(err)
		return comparison > 0
	})

	AssertEqualOrders(&suite.Suite, expected, input)
}

func (suite *supplyComparatorSuite) TestSortsByQuantityInDescending() {
	timeNow := time.Now()

	input := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(12), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(13), Timestamp: timeNow.UnixNano()},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(13), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(300), Qty: decimal.NewFromFloat(12), Timestamp: timeNow.UnixNano()},
	}

	sort.Slice(input, func(i, j int) bool {
		comparison, err := suite.comparator(input[i], input[j])
		suite.Require().NoError(err)
		return comparison > 0
	})

	AssertEqualOrders(&suite.Suite, expected, input)
}

func (suite *supplyComparatorSuite) TestSortsByTimestampInAscending() {
	timeNow := time.Now()

	input := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.Add(1 * time.Second).UnixNano()},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.Add(1 * time.Second).UnixNano()},
	}

	sort.Slice(input, func(i, j int) bool {
		comparison, err := suite.comparator(input[i], input[j])
		suite.Require().NoError(err)
		return comparison > 0
	})

	AssertEqualOrders(&suite.Suite, expected, input)
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
