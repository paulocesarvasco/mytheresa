package repository

import (
	"context"

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

func (r *ProductStore) ListProducts(ctx context.Context, limit int, offset int) ([]Product, int64, error) {
	var products []Product
	var total int64

	if err := r.db.
		WithContext(ctx).
		Model(&Product{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		WithContext(ctx).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Preload("Category").
		Preload("Variants").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
