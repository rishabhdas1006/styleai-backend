package models

type Order struct {
	ID         uint        `gorm:"primaryKey"`
	UserID     uint        `gorm:"index;not null"`
	Status     OrderStatus `gorm:"type:varchar(20);not null"`
	TotalPrice float64     `gorm:"not null"`

	Items []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)
