package categories

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/logs"
)

type CategoriesStore interface {
	ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]Category, int64, error)
	CreateCategory(ctx context.Context, code string, name string) (Category, error)
	CreateCategories(ctx context.Context, in []CreateCategoryInput) ([]Category, error)
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

func (s *Service) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]Category, int64, error) {
	return s.store.ListCategories(ctx, limit, offset, categoryCode)
}

func (s *Service) CreateCategory(ctx context.Context, code string, name string) (Category, error) {
	return s.store.CreateCategory(ctx, code, name)
}

func (s *Service) CreateCategories(ctx context.Context, in []CreateCategoryInput) ([]Category, error) {
	return nil, nil
}
