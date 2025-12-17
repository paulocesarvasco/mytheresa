package repository

import (
	"context"
	"errors"

	"github.com/mytheresa/go-hiring-challenge/internal/categories"
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
func (cs *CategoryStore) ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]categories.Category, int64, error) {

	var registers []Category
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

	if err := selectQuery.Find(&registers).Error; err != nil {
		cs.log.Error(ctx, "repository error fetching categories",
			"error", err)
		return nil, 0, errorsapi.ErrRepositoryFetchCategories
	}

	cats := make([]categories.Category, 0, len(registers))
	for _, r := range registers {
		cats = append(cats, categories.Category{
			Code: r.Code,
			Name: r.Name,
		})
	}

	return cats, total, nil
}

func (cs *CategoryStore) CreateCategory(ctx context.Context, code string, name string) (categories.Category, error) {
	category := Category{
		Code: code,
		Name: name,
	}

	err := cs.db.WithContext(ctx).
		Create(&category).
		Error

	if err != nil {
		cs.log.Error(ctx, "create category failed",
			"err", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return categories.Category{}, errorsapi.ErrRepositoryCategoryAlreadyExists
		}
		return categories.Category{}, errorsapi.ErrRepositoryCreateCategory
	}

	return categories.Category{Code: category.Code, Name: category.Name}, nil
}

func (cs *CategoryStore) CreateCategories(ctx context.Context, inputs []categories.CreateCategoryInput) ([]categories.Category, error) {
	return nil, nil

}
