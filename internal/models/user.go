package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:100;unique;not null" json:"username"`
	Password  string         `gorm:"not null" json:"password"` // Encriptada
	RoleID    uint           `json:"role_id"`
	Role      Role           `gorm:"constraint:OnDelete:SET NULL;" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
