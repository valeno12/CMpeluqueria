package dtos

type ClientDTO struct {
	Name     string `json:"name" example:"Valentino"`
	LastName string `json:"last_name" example:"Garcia Mendez"`
	Phone    string `json:"phone" example:"343534345"`
	Email    string `json:"email" example:"example@gmail.com"`
}

type GetClientDto struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name" example:"Valentino"`
	LastName     string                 `json:"last_name" example:"Garcia Mendez"`
	Phone        string                 `json:"phone" example:"343534345"`
	Email        string                 `json:"email" example:"example@gmail.com"`
	Appointments []ClientAppointmentDto `json:"appointments"`
}

type ClientAppointmentDto struct {
	ID              uint   `json:"id"`
	AppointmentDate string `json:"appointment_date" example:"30/09/2002 16:30"`
	Status          string `json:"status" example:"finalizado"`
}
