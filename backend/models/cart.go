package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Items     []Item         `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"items"`
	UserId    uint           `json:"user_id"`
	User      User           `json:"user"`
	Status    string         `gorm:"status;default:'pending'" json:"status"` // pending , confirmed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
