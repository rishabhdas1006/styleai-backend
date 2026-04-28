package models

type VariantImage struct {
	ID        uint   `gorm:"primaryKey"`
	VariantID string `gorm:"type:uuid;not null;index"`

	ImageURL string `gorm:"not null"`
	PublicID string `gorm:"not null"`

	Position int `gorm:"default:0"`

	Variant ProductVariant `gorm:"foreignKey:VariantID;references:ID;constraint:OnDelete:CASCADE"`
}
