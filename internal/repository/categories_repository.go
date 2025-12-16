package repository

import (
	"context"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"gorm.io/gorm"
)

type CategoryStore struct {
	db  *gorm.DB
	log logs.ApiLogger
}

func NewCategoryStore(db *gorm.DB) *CategoryStore {
	return &CategoryStore{
		db:  db,
		log: logs.Logger(),
	}
}
func (cs *CategoryStore) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]Category, int64, error) {

	var categories []Category
	var total int64

	countQuery := cs.db.WithContext(ctx).Model(&Category{})

	if categoryCode != "" {
		countQuery = countQuery.Where("code = ?", categoryCode)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		cs.log.Error(ctx, "repository error counting categories",
			"error", err)
		return nil, 0, errorsapi.ErrRepositoryCountCategories
	}

	selectQuery := cs.db.WithContext(ctx).
		Order("categories.code ASC").
		Limit(limit).
		Offset(offset)

	if categoryCode != "" {
		selectQuery = selectQuery.Where("code = ?", categoryCode)
	}

	if err := selectQuery.Find(&categories).Error; err != nil {
		cs.log.Error(ctx, "repository error fetching categories",
			"error", err)
		return nil, 0, errorsapi.ErrRepositoryFetchCategories
	}

	return categories, total, nil
}
