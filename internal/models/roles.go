package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"size:100;unique;not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"` // Relación many-to-many
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;unique;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}

type RolePermission struct {
	RoleID       uint           `gorm:"primaryKey" json:"role_id"`
	Role         Role           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"role"`
	PermissionID uint           `gorm:"primaryKey" json:"permission_id"`
	Permission   Permission     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permission"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-" swag:"-"`
}
