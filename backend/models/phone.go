package models

import (
	"time"

	"gorm.io/gorm"
)

type Phone struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Price     float32        `gorm:"price;default:0" json:"price"`
	BrandName string         `gorm:"bland_name" json:"brand_name"`
	ModelName string         `gorm:"model_name" json:"model_name"`
	OS        string         `gorm:"os" json:"os"`
	Amount    int            `gorm:"amount;default:0" json:"amount"`
	Image     []byte         `gorm:"image;type:longblob" json:"image"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
