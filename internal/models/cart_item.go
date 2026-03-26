package models

import "time"

type CartItem struct {
	ID        uint    `gorm:"primaryKey"`
	CartID    uint    `gorm:"index;not null"`
	VariantID uint    `gorm:"index;not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"` // snapshot price
	CreatedAt time.Time
	UpdatedAt time.Time
}
