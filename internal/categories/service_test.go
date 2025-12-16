package categories

import (
	"testing"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	tests := []struct {
		name               string
		fakeCategories     []repository.Category
		fakeTotal          int64
		fakeError          error
		expectedCategories CategoryPage
		expectedError      error
	}{
		{
			name: "list categories succeeds",
			fakeCategories: []repository.Category{
				{ID: 1, Code: "FOO", Name: "foo"},
				{ID: 2, Code: "BAR", Name: "bar"},
			},
			fakeTotal: 2,
			expectedCategories: CategoryPage{Categories: []CategoryView{
				{Code: "FOO", Name: "foo"},
				{Code: "BAR", Name: "bar"},
			}, Total: 2},
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

			products, err := service.ListCategories(t.Context(), 10, 0, "FOO")

			assert.Equal(t, tt.expectedCategories, products)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
