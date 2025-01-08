package dtos

import "time"

type CreateAppointmentDto struct {
	ClientID        uint   `json:"client_id" example:"1"`                       // ID del cliente
	AppointmentDate string `json:"appointment_date" example:"15:30 12/01/2025"` // Formato: HH:MM DD/MM/YYYY
	ServiceIds      []uint `json:"service_id" example:"1,2"`                    // IDs de los servicios asociados
}

type AppointmentServiceDto struct {
	ServiceID            uint    `json:"service_id" example:"1"`
	ServiceName          string  `json:"service_name" example:"Corte de cabello"`
	Price                float64 `json:"price" example:"1500.00"`
	EstimatedTimeMinutes uint    `json:"estimated_time_minutes" example:"30"`
}

type AppointmentByIDDto struct {
	ID              uint                    `json:"id" example:"1"`
	ClientID        uint                    `json:"client_id" example:"1"`
	ClientName      string                  `json:"client_name" example:"Juan Pérez"`
	Status          string                  `json:"status" example:"pendiente"`
	AppointmentDate string                  `json:"appointment_date" example:"12/01/2025 15:30"`
	Services        []AppointmentServiceDto `json:"services"`
	Products        []AppointmentProductDto `json:"products"`
	CreatedAt       time.Time               `json:"created_at" example:"2025-01-08T10:00:00Z"`
	UpdatedAt       time.Time               `json:"updated_at" example:"2025-01-08T12:00:00Z"`
}

type AllAppointmentDto struct {
	ID                   uint   `json:"id" example:"1"`
	ClientID             uint   `json:"client_id" example:"1"`
	ClientName           string `json:"client_name" example:"Juan Pérez"`
	Status               string `json:"status" example:"pendiente"`
	AppointmentDate      string `json:"appointment_date" example:"12/01/2025 15:30"`
	EstimatedTimeMinutes uint   `json:"estimated_time_minutes" example:"60"`
}

type FinalizeAppointmentDto struct {
	PaymentMethod string                          `json:"payment_method" example:"tarjeta"`
	Products      []FinalizeAppointmentProductDto `json:"products"`
}

type FinalizeAppointmentProductDto struct {
	ProductID uint    `json:"product_id" example:"1"`
	Quantity  float64 `json:"quantity" example:"2"`
}

type AppointmentProductDto struct {
	ProductID uint    `json:"product_id" example:"1"`
	Name      string  `json:"name" example:"Gel fijador"`
	Quantity  float64 `json:"quantity" example:"1"`
	Unit      string  `json:"unit" example:"unidad"` // Unidad del producto
}

type UpdateAppointmentProductsDto struct {
	Products []FinalizeAppointmentProductDto `json:"products"`
}
