package dtos

type UserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
}

type GetUserDto struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
}
