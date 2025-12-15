package catalogapi

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	"github.com/shopspring/decimal"
)

type fakeService struct {
	listProductsResp catalog.ProductPage
	listProductsErr  error

	detailProductResp catalog.ProductView
	detailProductErr  error
}

func NewFakeService() *fakeService {
	return &fakeService{}
}

func (f *fakeService) SetListProductsResponse(products catalog.ProductPage, err error) {
	f.listProductsResp = products
	f.listProductsErr = err
}

func (f *fakeService) SetDetailProductResponse(details catalog.ProductView, err error) {
	f.detailProductResp = details
	f.detailProductErr = err
}

func (f *fakeService) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) (catalog.ProductPage, error) {
	return f.listProductsResp, f.listProductsErr
}

func (f *fakeService) DetailProduct(ctx context.Context, code string) (catalog.ProductView, error) {
	return f.detailProductResp, f.detailProductErr
}
