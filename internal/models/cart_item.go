package models

import "time"

type CartItem struct {
	ID        uint   `gorm:"primaryKey"`
	CartID    uint   `gorm:"index;not null"`
	VariantID string `gorm:"type:uuid;index;not null"`

	Quantity int     `gorm:"not null"`
	Price    float64 `gorm:"not null"`

	Variant ProductVariant `gorm:"foreignKey:VariantID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
