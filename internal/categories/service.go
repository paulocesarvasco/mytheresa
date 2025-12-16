package categories

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

type CategoriesStore interface {
	ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]repository.Category, int64, error)
}

type Service struct {
	log   logs.ApiLogger
	store CategoriesStore
}

func New(store CategoriesStore) *Service {
	return &Service{
		log:   logs.Logger(),
		store: store}
}

func (s *Service) ListCategories(ctx context.Context, limit, offset int, categoryCode string) (CategoryPage, error) {
	res, total, err := s.store.ListCategories(ctx, limit, offset, categoryCode)
	if err != nil {
		return CategoryPage{}, err
	}

	categories := make([]CategoryView, len(res))
	for i, c := range res {
		categories[i] = CategoryView{
			Code: c.Code,
			Name: c.Name,
		}
	}

	return CategoryPage{Categories: categories, Total: total}, nil
}

func (s *Service) CreateCategory(ctx context.Context, code string, name string) (CategoryView, error) {
	s.log.Debug(ctx, "create category", "code", code, "name", name)
	return CategoryView{}, nil
}
