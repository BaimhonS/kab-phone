package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	TrackingNumber string         `gorm:"tracking_number;unique;not null" json:"tracking_number"`
	CartID         uint           `json:"cart_id"`
	Cart           Cart           `json:"cart"`
	TotalPrice     float32        `json:"total_price"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
