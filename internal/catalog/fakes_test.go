package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/shopspring/decimal"
)

type fakeStore struct {
	listProductsResp []repository.Product
	listTotal        int64
	listProductsErr  error

	getGetByCodeResp *repository.Product
	getGetByCodeErr  error
}

func NewFakeStore() *fakeStore {
	return &fakeStore{}
}

func (f *fakeStore) SetListProductsResponse(products []repository.Product, total int64, err error) {
	f.listProductsResp = products
	f.listTotal = total
	f.listProductsErr = err
}

func (f *fakeStore) SetGetByCodeResponse(details *repository.Product, err error) {
	f.getGetByCodeResp = details
	f.getGetByCodeErr = err
}

func (f *fakeStore) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]repository.Product, int64, error) {
	return f.listProductsResp, f.listTotal, f.listProductsErr
}

func (f *fakeStore) GetByCode(ctx context.Context, code string) (*repository.Product, error) {
	return f.getGetByCodeResp, f.getGetByCodeErr
}
