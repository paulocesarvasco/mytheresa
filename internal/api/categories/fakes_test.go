package categoriesapi

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/categories"
)

type fakeService struct {
	listCategoriesResp categories.CategoryPage
	listCategoriesErr  error
}

func NewFakeService() *fakeService {
	return &fakeService{}
}

func (f *fakeService) SetListCategoriesResponse(categories categories.CategoryPage, err error) {
	f.listCategoriesResp = categories
	f.listCategoriesErr = err
}

func (f *fakeService) ListCategories(ctx context.Context, limit, offset int, categoryCode string) (categories.CategoryPage, error) {
	return f.listCategoriesResp, f.listCategoriesErr
}
