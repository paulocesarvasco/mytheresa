package categoriesapi

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/categories"
)

type fakeService struct {
	listCategoriesResp  []categories.Category
	listCategoriesTotal int64
	listCategoriesErr   error

	createCategoryErr error
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
	f.createCategoryErr = err
}

func (f *fakeService) SetCreateCategoriesResponse(category categories.Category, err error) {
	f.createCategoryErr = err
}

func (f *fakeService) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]categories.Category, int64, error) {
	return f.listCategoriesResp, f.listCategoriesTotal, f.listCategoriesErr
}

func (f *fakeService) CreateCategory(ctx context.Context, code string, name string) error {
	return f.createCategoryErr
}

func (f *fakeService) CreateCategories(ctx context.Context, categories []categories.CreateCategoryInput) error {
	return f.createCategoryErr
}
