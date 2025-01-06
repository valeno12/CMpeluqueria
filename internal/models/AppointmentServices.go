package models

import (
	"time"

	"gorm.io/gorm"
)

type AppointmentService struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppointmentID uint           `gorm:"not null" json:"appointment_id"`
	Appointment   Appointment    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"appointment"`
	ServiceID     uint           `gorm:"not null" json:"service_id"`
	Service       Service        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"service"`
	Price         float64        `gorm:"not null" json:"price"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
