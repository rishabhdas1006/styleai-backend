package models

type ProductVariant struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint `gorm:"uniqueIndex:idx_variant"`

	Size  string `gorm:"uniqueIndex:idx_variant;not null"`
	Color string `gorm:"uniqueIndex:idx_variant;not null"`

	SKU string `gorm:"uniqueIndex;not null"`

	Price float64 `gorm:"not null"`
	Stock int     `gorm:"not null"`

	Product Product `gorm:"constraint:OnDelete:CASCADE"`
}
