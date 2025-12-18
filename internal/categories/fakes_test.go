package categories

import (
	"context"
)

type fakeStore struct {
	listCategoriesResp []Category
	listTotal          int64
	listProductsErr    error

	createCategoryErr error
}

func NewFakeStore() *fakeStore {
	return &fakeStore{}
}

func (f *fakeStore) SetListCategoriesResponse(categories []Category, total int64, err error) {
	f.listCategoriesResp = categories
	f.listTotal = total
	f.listProductsErr = err
}

func (f *fakeStore) SetCreateCategoryResponse(err error) {
	f.createCategoryErr = err
}

func (f *fakeStore) SetCreateCategoriesResponse(err error) {
	f.createCategoryErr = err
}

func (f *fakeStore) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]Category, int64, error) {
	return f.listCategoriesResp, f.listTotal, f.listProductsErr
}

func (f *fakeStore) CreateCategory(ctx context.Context, code string, name string) error {
	return f.createCategoryErr
}

func (f *fakeStore) CreateCategories(ctx context.Context, inputs []CreateCategoryInput) error {
	return f.createCategoryErr
}
