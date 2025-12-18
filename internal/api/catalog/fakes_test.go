package catalogapi

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	"github.com/shopspring/decimal"
)

type fakeService struct {
	listProductsResp  []catalog.Product
	listProductsTotal int64
	listProductsErr   error

	detailProductResp catalog.Product
	detailProductErr  error
}

func NewFakeService() *fakeService {
	return &fakeService{}
}

func (f *fakeService) SetListProductsResponse(products []catalog.Product, total int64, err error) {
	f.listProductsResp = products
	f.listProductsTotal = total
	f.listProductsErr = err
}

func (f *fakeService) SetDetailProductResponse(details catalog.Product, err error) {
	f.detailProductResp = details
	f.detailProductErr = err
}

func (f *fakeService) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]catalog.Product, int64, error) {
	return f.listProductsResp, f.listProductsTotal, f.listProductsErr
}

func (f *fakeService) DetailProduct(ctx context.Context, code string) (catalog.Product, error) {
	return f.detailProductResp, f.detailProductErr
}
