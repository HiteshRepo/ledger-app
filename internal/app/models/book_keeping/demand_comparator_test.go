package book_keeping_test

import (
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"sort"
	"testing"
	"time"
)

type demandComparatorSuite struct {
	suite.Suite
	comparator book_keeping.Comparator
}

func TestDemandComparatorSuite(t *testing.T) {
	suite.Run(t, new(demandComparatorSuite))
}

func (suite *demandComparatorSuite) SetupTest() {
	suite.comparator = book_keeping.ProvideDemandComparator()
}

func (suite *demandComparatorSuite) TestSortsByPriceInDescending() {
	timeNow := time.Now()

	input := []*order.Order{
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(100), Qty: decimal.NewFromFloat(10), Timestamp: timeNow.UnixNano()},
	}

	sort.Slice(input, func(i, j int) bool {
		comparison, err := suite.comparator(input[i], input[j])
		suite.Require().NoError(err)
		return comparison > 0
	})

	AssertEqualOrders(&suite.Suite, expected, input)
}

func (suite *demandComparatorSuite) TestSortsByQuantityInDescending() {
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

func (suite *demandComparatorSuite) TestSortsByTimestampInAscending() {
	timeNow := time.Now()

	input := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.Add(1*time.Second).UnixNano()},
	}

	expected := []*order.Order{
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.UnixNano()},
		{Price: decimal.NewFromFloat(200), Qty: decimal.NewFromFloat(11), Timestamp: timeNow.Add(1*time.Second).UnixNano()},
	}

	sort.Slice(input, func(i, j int) bool {
		comparison, err := suite.comparator(input[i], input[j])
		suite.Require().NoError(err)
		return comparison > 0
	})

	AssertEqualOrders(&suite.Suite, expected, input)
}