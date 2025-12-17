package categories

import (
	"context"
)

type fakeStore struct {
	listCategoriesResp []Category
	listTotal          int64
	listProductsErr    error

	createCategoryResp Category
	createCategoryErr  error
}

func NewFakeStore() *fakeStore {
	return &fakeStore{}
}

func (f *fakeStore) SetListCategoriesResponse(categories []Category, total int64, err error) {
	f.listCategoriesResp = categories
	f.listTotal = total
	f.listProductsErr = err
}

func (f *fakeStore) SetCreateCategoryResponse(category Category, err error) {
	f.createCategoryResp = category
	f.createCategoryErr = err
}

func (f *fakeStore) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]Category, int64, error) {
	return f.listCategoriesResp, f.listTotal, f.listProductsErr
}

func (f *fakeStore) CreateCategory(ctx context.Context, code string, name string) (Category, error) {
	return f.createCategoryResp, f.createCategoryErr
}

func (f *fakeStore) CreateCategories(ctx context.Context, in []CreateCategoryInput) ([]Category, error) {
	return nil, nil
}
