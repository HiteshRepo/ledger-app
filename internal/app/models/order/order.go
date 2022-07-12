package order

import (
	shopspring "github.com/shopspring/decimal"
)

type Order struct {
	Id        string
	Price     shopspring.Decimal
	Qty       shopspring.Decimal
	OrderType string
	Status    string
	Timestamp int64
}
