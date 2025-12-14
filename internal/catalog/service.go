package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

type ProductStore interface {
	ListProducts(ctx context.Context, limit int, offset int) ([]repository.Product, int64, error)
}

type ProductPage struct {
	Products []ProductView `json:"products"`
	Total    int64         `json:"total"`
}

type ProductView struct {
	CategoryName string  `json:"category_name"`
	CategoryCode string  `json:"category_code"`
	Code         string  `json:"code"`
	Price        float64 `json:"price"`
}

type Service struct {
	store ProductStore
}

func New(store ProductStore) *Service {
	return &Service{store: store}
}

func (s *Service) ListProducts(ctx context.Context, limit int, offset int) (ProductPage, error) {
	res, total, err := s.store.ListProducts(ctx, limit, offset)
	if err != nil {
		// TODO: improve error handler
		return ProductPage{}, err
	}

	products := make([]ProductView, len(res))
	for i, p := range res {
		products[i] = ProductView{
			CategoryName: p.Category.Name,
			CategoryCode: p.Category.Code,
			Code:         p.Code,
			Price:        p.Price.InexactFloat64(),
		}
	}

	// TODO: handle cases without results
	return ProductPage{Products: products, Total: total}, nil
}
