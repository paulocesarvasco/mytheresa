package catalogapi

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	"github.com/shopspring/decimal"
)

type fakeService struct {
	products catalog.ProductPage
	err      error
}

func NewFakeService(products catalog.ProductPage, err error) Service {
	return &fakeService{
		products: products,
		err:      err,
	}
}

func (f *fakeService) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) (catalog.ProductPage, error) {
	return f.products, f.err
}
