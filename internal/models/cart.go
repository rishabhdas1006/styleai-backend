package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"uniqueIndex;not null"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
