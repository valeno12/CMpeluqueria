package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ClientID        uint           `gorm:"not null" json:"client_id"`
	Client          Client         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"client"`
	Status          string         `gorm:"size:50;not null" json:"status"` // Ej: "pendiente", "cancelado", "finalizado"
	AppointmentDate time.Time      `gorm:"not null" json:"appointment_date"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
