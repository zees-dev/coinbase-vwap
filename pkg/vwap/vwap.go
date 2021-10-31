package vwap

import "github.com/shopspring/decimal"

type Vwapper interface {
	VWAP() decimal.Decimal
	Update(price, volume string) error
}
