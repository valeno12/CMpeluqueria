package dtos

import "time"

type CreateAppointmentDto struct {
	ClientID        uint   `json:"client_id"`
	AppointmentDate string `json:"appointment_date"` // Formato: HH:MM DD/MM/YYYY
	ServiceIds      []uint `json:"service_id"`
}

type AppointmentServiceDto struct {
	ServiceID            uint    `json:"service_id"`
	ServiceName          string  `json:"service_name"`
	Price                float64 `json:"price"`
	EstimatedTimeMinutes uint    `json:"estimated_time_minutes"`
}

type AppointmentByIDDto struct {
	ID              uint                    `json:"id"`
	ClientID        uint                    `json:"client_id"`
	ClientName      string                  `json:"client_name"`
	Status          string                  `json:"status"`
	AppointmentDate time.Time               `json:"appointment_date"`
	Services        []AppointmentServiceDto `json:"services"`
	Products        []AppointmentProductDto `json:"products"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type AllAppointmentDto struct {
	ID                   uint      `json:"id"`
	ClientID             uint      `json:"client_id"`
	ClientName           string    `json:"client_name"`
	Status               string    `json:"status"`
	AppointmentDate      time.Time `json:"appointment_date"`
	EstimatedTimeMinutes uint      `json:"estimated_time_minutes"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type FinalizeAppointmentDto struct {
	PaymentMethod string                          `json:"payment_method"`
	Products      []FinalizeAppointmentProductDto `json:"products"`
}

type FinalizeAppointmentProductDto struct {
	ProductID uint    `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

type AppointmentProductDto struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"` // Unidad del producto
}

type UpdateAppointmentProductsDto struct {
	Products []FinalizeAppointmentProductDto `json:"products"`
}
