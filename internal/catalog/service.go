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

func (s *Service) DetailProduct(ctx context.Context, code string) (ProductView, error) {
	product, err := s.store.GetByCode(ctx, code)
	if err != nil {
		return ProductView{}, err
	}

	variants := make([]VariantView, len(product.Variants))
	for i, v := range product.Variants {
		variants[i] = toVariantView(v, product.Price.InexactFloat64())
	}

	return ProductView{
		Category: product.Category.Code,
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Variants: variants,
	}, nil
}
