package models

type OrderItem struct {
	ID        uint   `gorm:"primaryKey"`
	OrderID   uint   `gorm:"index;not null"`
	VariantID string `gorm:"type:uuid;index;not null"`

	Quantity int     `gorm:"not null"`
	Price    float64 `gorm:"not null"`
}
