package categories

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

type fakeStore struct {
	listCategoriesResp []repository.Category
	listTotal          int64
	listProductsErr    error

	createCategoryResp repository.Category
	createCategoryErr  error
}

func NewFakeStore() *fakeStore {
	return &fakeStore{}
}

func (f *fakeStore) SetListCategoriesResponse(categories []repository.Category, total int64, err error) {
	f.listCategoriesResp = categories
	f.listTotal = total
	f.listProductsErr = err
}

func (f *fakeStore) SetCreateCategoryResponse(category repository.Category, err error) {
	f.createCategoryResp = category
	f.createCategoryErr = err
}

func (f *fakeStore) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]repository.Category, int64, error) {
	return f.listCategoriesResp, f.listTotal, f.listProductsErr
}

func (f *fakeStore) CreateCategory(ctx context.Context, code string, name string) (repository.Category, error) {
	return f.createCategoryResp, f.createCategoryErr
}
