package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Role      string `gorm:"type:varchar(20);default:user;index"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
}
