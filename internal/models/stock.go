package models

import "time"

type StockMovement struct {
	ID             uint      `gorm:"primaryKey"`
	ProductID      uint      `gorm:"not null"`
	Product        Product   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity       float64   `gorm:"not null"`                // Positivo para entrada, negativo para salida
	PackageCount   *float64  `json:"package_count,omitempty"` // Paquetes (NULL en salidas)
	ProductUnit    string    `gorm:"size:100" json:"product_unit"`
	UnitPerPackage *float64  `json:"unit_per_package,omitempty"` // Unidades por paquete (NULL en salidas)
	UnityPrice     *float64  `json:"unity_price,omitempty"`      // Precio unitario (NULL en salidas)
	Reason         string    `gorm:"size:255"`
	CreatedAt      time.Time `json:"created_at"`
}
