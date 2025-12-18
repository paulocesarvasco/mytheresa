package repository

// Category is the database model for product categories.
type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`

	Products []Product `gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "categories"
}
