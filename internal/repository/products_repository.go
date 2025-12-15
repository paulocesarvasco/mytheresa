package repository

import (
	"context"
	"errors"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductStore struct {
	db  *gorm.DB
	log logs.ApiLogger
}

func New(db *gorm.DB) *ProductStore {
	return &ProductStore{
		db:  db,
		log: logs.Logger(),
	}
}
func (ps *ProductStore) ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) ([]Product, int64, error) {

	var products []Product
	var total int64

	countQuery := ps.db.WithContext(ctx).Model(&Product{})

	if categoryCode != "" {
		countQuery = countQuery.
			Joins("Category").
			Where(`"Category"."code" = ?`, categoryCode)
	}
	if maxPrice != nil {
		countQuery = countQuery.Where("products.price < ?", maxPrice)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		ps.log.Error(ctx, "repository error counting products", "error", err)
		return nil, 0, errorsapi.ErrRepositoryCountProducts
	}

	selectQuery := ps.db.WithContext(ctx).
		Order("products.price DESC").
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
		ps.log.Error(ctx, "repository error fetching products", "error", err)
		return nil, 0, errorsapi.ErrRepositoryFetchProducts
	}

	return products, total, nil
}

func (ps *ProductStore) GetByCode(ctx context.Context, code string) (*Product, error) {
	var product Product

	err := ps.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Where("code = ?", code).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsapi.ErrProductNotFound
		}
		ps.log.Error(ctx, "repository error fetching product details", "error", err)
		return nil, errorsapi.ErrRepositoryFetchProduct
	}

	return &product, nil
}
