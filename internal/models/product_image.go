package models

type ProductImage struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	ImageURL  string `gorm:"not null"`

	Product Product `gorm:"constraint:OnDelete:CASCADE"`
}
