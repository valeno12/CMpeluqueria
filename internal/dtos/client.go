package dtos

type ClientDTO struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type GetClientDto struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}
