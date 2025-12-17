package categoriesapi

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/categories"
)

type fakeService struct {
	listCategoriesResp  []categories.Category
	listCategoriesTotal int64
	listCategoriesErr   error

	createCategoryResp categories.Category
	createCategoryErr  error
}

func NewFakeService() *fakeService {
	return &fakeService{}
}

func (f *fakeService) SetListCategoriesResponse(categories []categories.Category, total int64, err error) {
	f.listCategoriesResp = categories
	f.listCategoriesTotal = total
	f.listCategoriesErr = err
}

func (f *fakeService) SetCreateCategoryResponse(category categories.Category, err error) {
	f.createCategoryResp = category
	f.createCategoryErr = err
}

func (f *fakeService) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]categories.Category, int64, error) {
	return f.listCategoriesResp, f.listCategoriesTotal, f.listCategoriesErr
}

func (f *fakeService) CreateCategory(ctx context.Context, code string, name string) (categories.Category, error) {
	return f.createCategoryResp, f.createCategoryErr
}

func (f *fakeService) CreateCategories(ctx context.Context, categories []categories.CreateCategoryInput) ([]categories.Category, error) {
	return nil, nil
}
