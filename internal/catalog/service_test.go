package catalog

import (
	"testing"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	tests := []struct {
		name             string
		fakeProducts     []repository.Product
		fakeTotal        int64
		fakeError        error
		expectedProducts ProductPage
		expectedError    error
	}{
		{
			name:             "list products succeeds",
			fakeProducts:     []repository.Product{{ID: 1, Code: "PROD001", Price: decimal.New(10, 0), CategoryID: 1, Category: repository.Category{Code: "FOO"}}},
			fakeTotal:        1,
			expectedProducts: ProductPage{Products: []ProductView{{Category: "FOO", Code: "PROD001", Price: 10}}, Total: 1},
		},
		{
			name:          "list products fails on repository error",
			fakeError:     errorsapi.ErrRepositoryCountProducts,
			expectedError: errorsapi.ErrRepositoryCountProducts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewFakeStore()
			store.SetListProductsResponse(tt.fakeProducts, tt.fakeTotal, tt.fakeError)

			service := New(store)

			price := decimal.New(10, 0)
			products, err := service.ListProducts(t.Context(), 10, 0, "FOO", &price)

			assert.Equal(t, tt.expectedProducts, products)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestDetailProduct(t *testing.T) {
	tests := []struct {
		name            string
		fakeDetails     repository.Product
		fakeError       error
		expectedDetails ProductView
		expectedError   error
	}{
		{
			name: "detail product uses variant price when present",
			fakeDetails: repository.Product{ID: 1, Code: "PROD001", Price: *priceGenerator("9.99"), CategoryID: 1, Category: repository.Category{Code: "FOO"},
				Variants: []repository.Variant{{ID: 1, ProductID: 1, Name: "Variant A", SKU: "SKU001A"}}},
			expectedDetails: ProductView{Category: "FOO", Code: "PROD001", Price: 9.99,
				Variants: []VariantView{{Name: "Variant A", SKU: "SKU001A", Price: 9.99}}},
		},
		{
			name: "detail product applies price inheritance when variant price is not set",
			fakeDetails: repository.Product{ID: 1, Code: "PROD001", Price: *priceGenerator("9.99"), CategoryID: 1, Category: repository.Category{Code: "FOO"},
				Variants: []repository.Variant{{ID: 1, ProductID: 1, Name: "Variant B", SKU: "SKU001B", Price: priceGenerator("7.99")}}},
			expectedDetails: ProductView{Category: "FOO", Code: "PROD001", Price: 9.99,
				Variants: []VariantView{{Name: "Variant B", SKU: "SKU001B", Price: 7.99}}},
		},
		{
			name:          "detail product returns error on repository failure",
			fakeError:     errorsapi.ErrRepositoryFetchProduct,
			expectedError: errorsapi.ErrRepositoryFetchProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewFakeStore()
			store.SetGetByCodeResponse(&tt.fakeDetails, tt.fakeError)

			service := New(store)

			details, err := service.DetailProduct(t.Context(), "FOO")

			assert.Equal(t, tt.expectedDetails, details)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
