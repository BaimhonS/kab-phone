package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primarykey;autoIncrement" json:"id"`
	Username    string         `gorm:"username;unique;not null" json:"username"`
	FirstName   string         `gorm:"first_name;not null" json:"first_name"`
	LastName    string         `gorm:"last_name;not null" json:"last_name"`
	PhoneNumber string         `gorm:"phone_number;not null" json:"phone_number"`
	Password    string         `gorm:"password;not null" json:"password"`
	LineID      string         `gorm:"line_id" json:"line_id"`
	Address     string         `gorm:"address" json:"address"`
	Age         int            `gorm:"age" json:"age"`
	BirthDate   time.Time      `gorm:"birth_date" json:"birth_date"`
	Role        string         `gorm:"role;default:'guess'" json:"role"` // guess , admin
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
