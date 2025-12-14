package repository

// Category is the database model for product categories.
type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"uniqueIndex;not null"`
}

func (Category) TableName() string {
	return "categories"
}
