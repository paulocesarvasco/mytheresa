package categories

import (
	"testing"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestListProducts(t *testing.T) {
	tests := []struct {
		name               string
		fakeCategories     []Category
		fakeTotal          int64
		fakeError          error
		expectedCategories []Category
		expectedTotal      int64
		expectedError      error
	}{
		{
			name: "list categories succeeds",
			fakeCategories: []Category{
				{Code: "FOO", Name: "foo"},
				{Code: "BAR", Name: "bar"},
			},
			fakeTotal: 2,
			expectedCategories: []Category{
				{Code: "FOO", Name: "foo"},
				{Code: "BAR", Name: "bar"},
			},
			expectedTotal: 2,
		},
		{
			name:          "list categories fails on repository error",
			fakeError:     errorsapi.ErrRepositoryFetchCategories,
			expectedError: errorsapi.ErrRepositoryFetchCategories,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewFakeStore()
			store.SetListCategoriesResponse(tt.fakeCategories, tt.fakeTotal, tt.fakeError)

			service := New(store)

			products, total, err := service.ListCategories(t.Context(), 10, 0, "FOO")

			assert.Equal(t, tt.expectedCategories, products)
			assert.Equal(t, tt.expectedTotal, total)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name          string
		categoryName  string
		categoryCode  string
		fakeCategory  Category
		fakeError     error
		expectedError error
	}{
		{
			name:         "create category succeeds",
			fakeCategory: Category{Code: "FOO", Name: "foo"},
		},
		{
			name:          "create category fails on repository error",
			fakeError:     errorsapi.ErrRepositoryCreateCategory,
			expectedError: errorsapi.ErrRepositoryCreateCategory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewFakeStore()
			store.SetCreateCategoryResponse(tt.fakeError)

			service := New(store)
			err := service.CreateCategory(t.Context(), tt.categoryCode, tt.categoryName)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
