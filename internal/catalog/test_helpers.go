package catalog

import "github.com/shopspring/decimal"

func priceGenerator(s string) *decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return &d
}
