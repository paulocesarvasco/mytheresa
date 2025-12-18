package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/shopspring/decimal"
)

func mountVariant(v repository.Variant, productPrice *decimal.Decimal) Variant {
	if v.Price == nil {
		return Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: productPrice,
		}
	}

	return Variant{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: v.Price,
	}
}
