package catalog

import "github.com/mytheresa/go-hiring-challenge/internal/repository"

func toVariantView(v repository.Variant, productPrice float64) VariantView {
	if v.Price == nil {
		return VariantView{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: productPrice,
		}
	}

	return VariantView{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: v.Price.InexactFloat64(),
	}
}
