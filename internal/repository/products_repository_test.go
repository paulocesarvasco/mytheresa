package repository

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mytheresa/go-hiring-challenge/internal/database/testutil"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
)

func TestListProducts(t *testing.T) {
	tests := []struct {
		name             string
		limit            int
		offset           int
		categoryCode     string
		maxPrice         *decimal.Decimal
		timeout          time.Duration
		expectedProducts []Product
		expectedTotal    int64
		expectedError    error
	}{
		{
			name:    "retrieve complete products list",
			limit:   10,
			timeout: 60 * time.Second,
			expectedProducts: []Product{
				{ID: 5, Code: "PROD005", Price: decimal.New(2299, -2), CategoryID: 3},
				{ID: 7, Code: "PROD007", Price: decimal.New(1820, -2), CategoryID: 1},
				{ID: 4, Code: "PROD004", Price: decimal.New(1500, -2), CategoryID: 1},
				{ID: 2, Code: "PROD002", Price: decimal.New(1249, -2), CategoryID: 2},
				{ID: 1, Code: "PROD001", Price: decimal.New(1099, -2), CategoryID: 1},
				{ID: 8, Code: "PROD008", Price: decimal.New(999, -2), CategoryID: 3},
				{ID: 3, Code: "PROD003", Price: decimal.New(875, -2), CategoryID: 3},
				{ID: 6, Code: "PROD006", Price: decimal.New(550, -2), CategoryID: 2},
			},
			expectedTotal: 8,
			expectedError: nil,
		},
		{
			name:    "get 1st product from list",
			limit:   1,
			timeout: 60 * time.Second,
			expectedProducts: []Product{
				{ID: 5, Code: "PROD005", Price: decimal.New(2299, -2), CategoryID: 3},
			},
			expectedTotal: 8,
			expectedError: nil,
		},
		{
			name:    "get last product from list",
			limit:   1,
			offset:  7,
			timeout: 60 * time.Second,
			expectedProducts: []Product{
				{ID: 6, Code: "PROD006", Price: decimal.New(550, -2), CategoryID: 2},
			},
			expectedTotal: 8,
			expectedError: nil,
		},
		{
			name:         "search by category",
			limit:        10,
			timeout:      60 * time.Second,
			categoryCode: "CLOTHING",
			expectedProducts: []Product{
				{ID: 7, Code: "PROD007", Price: decimal.New(1820, -2), CategoryID: 1},
				{ID: 4, Code: "PROD004", Price: decimal.New(1500, -2), CategoryID: 1},
				{ID: 1, Code: "PROD001", Price: decimal.New(1099, -2), CategoryID: 1},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:         "search by category with limit without offset",
			limit:        1,
			timeout:      60 * time.Second,
			categoryCode: "CLOTHING",
			expectedProducts: []Product{
				{ID: 7, Code: "PROD007", Price: decimal.New(1820, -2), CategoryID: 1},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:         "search by category with limit and offset",
			limit:        1,
			offset:       1,
			timeout:      60 * time.Second,
			categoryCode: "CLOTHING",
			expectedProducts: []Product{
				{ID: 4, Code: "PROD004", Price: decimal.New(1500, -2), CategoryID: 1},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:         "search by category with limit and offset and max_price",
			limit:        1,
			offset:       1,
			maxPrice:     priceGenerator("18.00"),
			timeout:      60 * time.Second,
			categoryCode: "CLOTHING",
			expectedProducts: []Product{
				{ID: 1, Code: "PROD001", Price: decimal.New(1099, -2), CategoryID: 1},
			},
			expectedTotal: 2,
			expectedError: nil,
		},
		{
			name:     "set max price limit",
			limit:    10,
			maxPrice: priceGenerator("10.00"),
			timeout:  60 * time.Second,
			expectedProducts: []Product{
				{ID: 8, Code: "PROD008", Price: decimal.New(999, -2), CategoryID: 3},
				{ID: 3, Code: "PROD003", Price: decimal.New(875, -2), CategoryID: 3},
				{ID: 6, Code: "PROD006", Price: decimal.New(550, -2), CategoryID: 2},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:          "database connection timeout",
			limit:         10,
			timeout:       0 * time.Second,
			expectedError: errorsapi.ErrRepositoryCountProducts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := testutil.StartPostgresContainer(t, context.Background())
			t.Cleanup(cleanup)

			store := NewProductStore(db)
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			products, total, err := store.ListProducts(ctx, tt.limit, tt.offset, tt.categoryCode, tt.maxPrice)
			assert.Equal(t, tt.expectedTotal, total)
			assert.Equal(t, tt.expectedError, err)

			for i, p := range products {
				assert.Equal(t, tt.expectedProducts[i].ID, p.ID)
				assert.Equal(t, tt.expectedProducts[i].Code, p.Code)
				assert.Equal(t, tt.expectedProducts[i].CategoryID, p.CategoryID)
				require.True(t, tt.expectedProducts[i].Price.Equal(p.Price))
			}

		})
	}

}

func TestGetByCode(t *testing.T) {
	tests := []struct {
		name            string
		productCode     string
		timeout         time.Duration
		expectedProduct *Product
		expectedError   error
	}{
		{
			name:        "get product detail succeeds",
			productCode: "PROD001",
			timeout:     60 * time.Second,
			expectedProduct: &Product{ID: 1, Code: "PROD001", Price: decimal.New(1099, -2),
				CategoryID: 1, Category: Category{ID: 1, Code: "CLOTHING", Name: "Clothing"},
				Variants: []Variant{
					{ID: 1, ProductID: 1, Name: "Variant A", SKU: "SKU001A", Price: priceGenerator("11.99")},
					{ID: 2, ProductID: 1, Name: "Variant B", SKU: "SKU001B"},
					{ID: 3, ProductID: 1, Name: "Variant C", SKU: "SKU001C"},
				}},
			expectedError: nil,
		},
		{
			name:            "product not found",
			productCode:     "PROD010",
			timeout:         60 * time.Second,
			expectedProduct: nil,
			expectedError:   errorsapi.ErrProductNotFound,
		},
		{
			name:            "database connection timeout",
			productCode:     "PROD001",
			timeout:         0 * time.Second,
			expectedProduct: nil,
			expectedError:   errorsapi.ErrRepositoryFetchProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := testutil.StartPostgresContainer(t, context.Background())
			t.Cleanup(cleanup)

			store := NewProductStore(db)
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			product, err := store.GetByCode(ctx, tt.productCode)

			assert.Equal(t, tt.expectedError, err)

			if product == nil {
				assert.Equal(t, tt.expectedProduct, product)
				return
			}

			assert.Equal(t, tt.expectedProduct.ID, product.ID)
			assert.Equal(t, tt.expectedProduct.Code, product.Code)
			require.True(t, tt.expectedProduct.Price.Equal(product.Price))
			assert.Equal(t, tt.expectedProduct.CategoryID, product.CategoryID)
			assert.Equal(t, tt.expectedProduct.Category, product.Category)

			for i, v := range product.Variants {
				assert.Equal(t, tt.expectedProduct.Variants[i].ID, v.ID)
				assert.Equal(t, tt.expectedProduct.Variants[i].Name, v.Name)
				assert.Equal(t, tt.expectedProduct.Variants[i].SKU, v.SKU)
				if v.Price == nil {
					assert.Equal(t, tt.expectedProduct.Variants[i].Price, v.Price)
				} else {
					require.True(t, tt.expectedProduct.Variants[i].Price.Equal(*v.Price))
				}
			}
		})
	}

}
