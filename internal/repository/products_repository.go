package repository

import (
	"context"

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
		return nil, 0, err
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
		return nil, 0, err
	}

	return products, total, nil
}
