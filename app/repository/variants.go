package repository

import (
	"github.com/shopspring/decimal"
)

// Variant is the database model for product variants.
type Variant struct {
	ID        uint            `gorm:"primaryKey"`
	ProductID uint            `gorm:"not null"`
	Name      string          `gorm:"not null"`
	SKU       string          `gorm:"uniqueIndex;not null"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

func (Variant) TableName() string {
	return "product_variants"
}
