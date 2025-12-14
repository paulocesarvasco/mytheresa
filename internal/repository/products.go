package repository

import (
	"github.com/shopspring/decimal"
)

// Product is the database model for catalog products.
type Product struct {
	ID    uint            `gorm:"primaryKey"`
	Code  string          `gorm:"uniqueIndex;not null"`
	Price decimal.Decimal `gorm:"type:decimal(10,2);not null"`

	CategoryID uint     `gorm:"not null;index"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID"`

	Variants []Variant `gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}
