package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Description   string         `gorm:"size:255" json:"description"`
	Price         float64        `gorm:"not null" json:"price"`
	EstimatedTime uint           `gorm:"not null" json:"estimated_time"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
