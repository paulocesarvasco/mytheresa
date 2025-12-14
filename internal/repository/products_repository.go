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

func (r *ProductStore) ListProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
