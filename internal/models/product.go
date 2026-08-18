package models

import "time"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Brand       string

	CategoryID uint   `gorm:"index"`
	Gender     string `gorm:"type:varchar(10);index;not null"`

	CreatedByID uint `gorm:"not null;index"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	MinPrice float64 `gorm:"not null;default:0"`

	PrimaryImageURL      string `gorm:"not null"`
	PrimaryImagePublicID string `gorm:"not null"`

	Category Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}
