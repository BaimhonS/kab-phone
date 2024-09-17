package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Amount    int            `gorm:"amount;default:0" json:"amount"`
	PhoneID   uint           `json:"phone_id"`
	Phone     Phone          `gorm:"constraint:OnDelete:CASCADE;" json:"phone"`
	CartID    uint           `json:"cart_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
