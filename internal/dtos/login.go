package dtos

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAnswerDto struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
