package repository

import (
	"context"
	"errors"
	"strings"

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

func (cs *CategoryStore) CreateCategory(ctx context.Context, code string, name string) error {
	category := Category{
		Code: code,
		Name: name,
	}

	err := cs.db.WithContext(ctx).
		Create(&category).
		Error

	if err != nil {
		cs.log.Error(ctx, "failed to insert register on database",
			"err", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			cs.log.Warn(ctx, "conflicting category codes already exist in database", "code", code)
			return errorsapi.ErrRepositoryCategoryAlreadyExists
		}
		return errorsapi.ErrRepositoryCreateCategory
	}

	return nil
}

func (cs *CategoryStore) CreateCategories(ctx context.Context, inputs []categories.CreateCategoryInput) error {
	if len(inputs) == 0 {
		return errorsapi.ErrRepositoryEmptyCategoriesInputList
	}

	newRegisters := make([]Category, len(inputs))
	codes := make([]string, len(inputs))

	for i, in := range inputs {
		code := strings.TrimSpace(in.Code)
		newRegisters[i] = Category{Code: code, Name: in.Name}
		codes[i] = code
	}

	err := cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&newRegisters).Error
	})

	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		cs.log.Warn(ctx, "batch category creation conflict", "codes", codes)
		var existing []string
		qerr := cs.db.WithContext(ctx).
			Model(&[]Category{}).
			Where("code IN ?", codes).
			Pluck("code", &existing).Error

		if qerr != nil {
			cs.log.Warn(ctx, "failed to lookup existing category codes after conflict", "err", qerr)
		} else if len(existing) > 0 {
			cs.log.Warn(ctx, "conflicting category codes already exist in database", "codes", existing)
		}

		return errorsapi.ErrRepositoryCategoryAlreadyExists
	} else if err != nil {
		cs.log.Error(ctx, "failed to insert registers on database", "err", err)
		return errorsapi.ErrRepositoryCreateCategory
	}
	return nil
}
