package dtos

type ServiceDto struct {
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	Price                float64 `json:"price"`
	EstimatedTimeMinutes uint    `json:"estimated_time_minutes"`
	EstimatedTimeHours   uint    `json:"estimated_time_hours"`
}
