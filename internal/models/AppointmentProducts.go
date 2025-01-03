package models

import (
	"time"

	"gorm.io/gorm"
)

type AppointmentProduct struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppointmentID uint           `gorm:"not null" json:"appointment_id"`
	Appointment   Appointment    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"appointment"`
	ProductID     uint           `gorm:"not null" json:"product_id"`
	Product       Product        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product"`
	Quantity      uint           `gorm:"not null" json:"quantity"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
