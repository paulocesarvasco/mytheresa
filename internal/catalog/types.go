package catalog

import "github.com/shopspring/decimal"

type Product struct {
	Category string
	Code     string
	Price    *decimal.Decimal

	Variants []Variant
}

type Variant struct {
	Name  string
	SKU   string
	Price *decimal.Decimal
}
