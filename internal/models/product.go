package models

import "time"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Brand       string
	CategoryID  uint

	CreatedAt time.Time
	UpdatedAt time.Time

	MinPrice float64 `gorm:"not null;default:0"`

	Category Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID"`
}
