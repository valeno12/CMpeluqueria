package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100,not null" json:"name"`
	LastName  string    `gorm:"size:100" json:"last_name"`
	Phone     string    `gorm:"size:15" json:"phone"`
	Email     string    `gorm:"size:100" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
