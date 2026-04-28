package models

type ProductVariant struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	ProductID uint   `gorm:"index;not null"`

	Size  string `gorm:"not null;uniqueIndex:idx_variant"`
	Color string `gorm:"not null;uniqueIndex:idx_variant"`

	SKU string `gorm:"uniqueIndex;not null"`

	Price float64 `gorm:"not null"`
	Stock int     `gorm:"not null"`

	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`

	Images []VariantImage `gorm:"foreignKey:VariantID"`
}
