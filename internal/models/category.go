package models

import "time"

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`

	CreatedByID uint `gorm:"not null;index"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time

	Products []Product `gorm:"foreignKey:CategoryID"`
}
