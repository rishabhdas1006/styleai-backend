package models

import "time"

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex;not null"`
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}
