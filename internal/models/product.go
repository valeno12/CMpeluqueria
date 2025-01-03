package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Unit            string         `gorm:"size:100;not null" json:"unit"` //ej ml, unidades
	Quantity        uint           `gorm:"not null" json:"quantity"`
	Low_stock_alert uint           `gorm:"not null" json:"low_stock_alert"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
