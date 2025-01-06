package models

import "time"

type StockMovement struct {
	ID         uint      `gorm:"primaryKey"`
	ProductID  uint      `gorm:"not null"`
	Product    Product   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity   float64   `gorm:"not null"` // Positivo para entradas, negativo para salidas
	Reason     string    `gorm:"size:255"`
	UnityPrice float64   `gorm:"not null"`
	TotalPrice float64   `gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}
