package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Role      string `gorm:"default:user"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
}
