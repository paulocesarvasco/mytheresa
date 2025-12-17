package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mytheresa/go-hiring-challenge/internal/categories"
	"github.com/mytheresa/go-hiring-challenge/internal/database/testutil"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
)

func TestListCategories(t *testing.T) {
	tests := []struct {
		name               string
		limit              int
		offset             int
		categoryCode       string
		timeout            time.Duration
		expectedCategories []categories.Category
		expectedTotal      int64
		expectedError      error
	}{
		{
			name:    "retrieve entire categories list",
			limit:   10,
			timeout: 60 * time.Second,
			expectedCategories: []categories.Category{
				{Code: "ACCESSORIES", Name: "Accessories"},
				{Code: "CLOTHING", Name: "Clothing"},
				{Code: "SHOES", Name: "Shoes"},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:    "get 1st category from list",
			limit:   1,
			timeout: 60 * time.Second,
			expectedCategories: []categories.Category{
				{Code: "ACCESSORIES", Name: "Accessories"},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:    "get last category from list",
			limit:   1,
			offset:  2,
			timeout: 60 * time.Second,
			expectedCategories: []categories.Category{
				{Code: "SHOES", Name: "Shoes"},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:    "retrieve 2 items from list",
			limit:   2,
			timeout: 60 * time.Second,
			expectedCategories: []categories.Category{
				{Code: "ACCESSORIES", Name: "Accessories"},
				{Code: "CLOTHING", Name: "Clothing"},
			},
			expectedTotal: 3,
			expectedError: nil,
		},
		{
			name:         "get category info",
			limit:        10,
			categoryCode: "ACCESSORIES",
			timeout:      60 * time.Second,
			expectedCategories: []categories.Category{
				{Code: "ACCESSORIES", Name: "Accessories"},
			},
			expectedTotal: 1,
			expectedError: nil,
		},
		{
			name:          "database connection timeout",
			limit:         10,
			timeout:       0 * time.Second,
			expectedError: errorsapi.ErrRepositoryCountCategories,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := testutil.StartPostgresContainer(t, context.Background())
			t.Cleanup(cleanup)

			store := NewCategoryStore(db)
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			categories, total, err := store.ListCategories(ctx, tt.limit, tt.offset, tt.categoryCode)
			assert.Equal(t, tt.expectedCategories, categories)
			assert.Equal(t, tt.expectedTotal, total)
			assert.Equal(t, tt.expectedError, err)

		})
	}

}

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name          string
		categoryCode  string
		categoryName  string
		timeout       time.Duration
		expectedError error
	}{
		{
			name:          "create category succeeds",
			categoryCode:  "FOO",
			categoryName:  "foo",
			timeout:       60 * time.Second,
			expectedError: nil,
		},
		{
			name:          "category already exists",
			categoryCode:  "SHOES",
			categoryName:  "shoes",
			timeout:       60 * time.Second,
			expectedError: errorsapi.ErrRepositoryCategoryAlreadyExists,
		},
		{
			name:          "database connection timeout",
			categoryCode:  "FOO",
			categoryName:  "foo",
			expectedError: errorsapi.ErrRepositoryCreateCategory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := testutil.StartPostgresContainer(t, context.Background())
			t.Cleanup(cleanup)

			store := NewCategoryStore(db)
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			err := store.CreateCategory(ctx, tt.categoryCode, tt.categoryName)
			assert.Equal(t, tt.expectedError, err)

		})
	}

}
