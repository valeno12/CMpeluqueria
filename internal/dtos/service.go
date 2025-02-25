package dtos

type ServiceDto struct {
	Name                 string  `json:"name" example:"Corte de pelo"`
	Description          string  `json:"description" example:"Corte de pelo clasico"`
	Price                float64 `json:"price" example:"10000"`
	EstimatedTimeMinutes uint    `json:"estimated_time_minutes" example:"30"`
	EstimatedTimeHours   uint    `json:"estimated_time_hours" example:"1"`
}

type GetServiceDto struct {
	ID            uint    `json:"id" example:"1"`
	Name          string  `json:"name" example:"Corte de pelo"`
	Description   string  `json:"description" example:"Corte de pelo clasico"`
	Price         float64 `json:"price" example:"10000"`
	EstimatedTime uint    `json:"estimated_time_minutes" example:"90"`
}
