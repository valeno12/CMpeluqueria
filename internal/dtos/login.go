package dtos

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAnswerDto struct {
	Username string `json:"username" example:"admin"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzYzODMzMzcsImlhdCI6MTczNjI5NjkzNywicm"`
}
