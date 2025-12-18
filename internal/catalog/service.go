package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/shopspring/decimal"
)

type ProductStore interface {
	ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]repository.Product, int64, error)
	GetByCode(ctx context.Context, code string) (*repository.Product, error)
}

type Service struct {
	log   logs.ApiLogger
	store ProductStore
}

func New(store ProductStore) *Service {
	return &Service{
		log:   logs.Logger(),
		store: store}
}

func (s *Service) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]Product, int64, error) {
	res, total, err := s.store.ListProducts(ctx, limit, offset, categoryCode, maxPrice)
	if err != nil {
		return nil, 0, err
	}
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Category: p.Category.Code,
			Code:     p.Code,
			Price:    &p.Price,
		}
	}
	return products, total, nil
}

func (s *Service) DetailProduct(ctx context.Context, code string) (Product, error) {
	product, err := s.store.GetByCode(ctx, code)
	if err != nil {
		return Product{}, err
	}

	variants := make([]Variant, len(product.Variants))
	for i, v := range product.Variants {
		variants[i] = mountVariant(v, &product.Price)
	}

	return Product{
		Category: product.Category.Code,
		Code:     product.Code,
		Price:    &product.Price,
		Variants: variants,
	}, nil
}
