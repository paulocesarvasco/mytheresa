package repository

import (
	"context"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductStore struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ProductStore {
	return &ProductStore{
		db: db,
	}
}
func (r *ProductStore) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]Product, int64, error) {
	log := logs.NewLogger()

	var products []Product
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&Product{})

	if categoryCode != "" {
		countQuery = countQuery.
			Joins("Category").
			Where(`"Category"."code" = ?`, categoryCode)
	}
	if maxPrice != nil {
		countQuery = countQuery.Where("products.price < ?", maxPrice)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		log.Error(ctx, "count", "error", err)
		return nil, 0, errorsapi.ErrRepositoryCountProducts
	}

	selectQuery := r.db.WithContext(ctx).
		Order("products.id ASC").
		Limit(limit).
		Offset(offset).
		Preload("Category").
		Preload("Variants")

	if categoryCode != "" {
		selectQuery = selectQuery.Joins("Category").
			Where(`"Category"."code" = ?`, categoryCode)
	}

	if maxPrice != nil {
		selectQuery = selectQuery.Where("products.price < ?", maxPrice)
	}

	if err := selectQuery.Find(&products).Error; err != nil {
		log.Error(ctx, "select", "error", err)
		return nil, 0, errorsapi.ErrRepositoryFetchProducts
	}

	return products, total, nil
}
