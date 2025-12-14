package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/shopspring/decimal"
)

type ProductStore interface {
	ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]repository.Product, int64, error)
}

type ProductPage struct {
	Products []ProductView `json:"products"`
	Total    int64         `json:"total"`
}

type ProductView struct {
	Category string  `json:"category_code"`
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
}

type Service struct {
	store ProductStore
}

func New(store ProductStore) *Service {
	return &Service{store: store}
}

func (s *Service) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) (ProductPage, error) {
	res, total, err := s.store.ListProducts(ctx, limit, offset, categoryCode, maxPrice)
	if err != nil {
		return ProductPage{}, err
	}

	products := make([]ProductView, len(res))
	for i, p := range res {
		products[i] = ProductView{
			Category: p.Category.Code,
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
		}
	}

	// TODO: handle cases without results
	return ProductPage{Products: products, Total: total}, nil
}
