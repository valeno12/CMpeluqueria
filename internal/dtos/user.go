package dtos

// UserDto representa los datos necesarios para crear un usuario
type UserDto struct {
	Username string `json:"username" example:"nuevo_usuario"` // Nombre único del usuario
	Password string `json:"password" example:"contraseña123"` // Contraseña del usuario
	RoleID   uint   `json:"role_id" example:"2"`              // ID del rol asignado al usuario
}

// GetUserDto representa los datos de salida al obtener un usuario
type GetUserDto struct {
	ID       uint   `json:"id" example:"1"`                    // ID del usuario
	Username string `json:"username" example:"admin"`          // Nombre de usuario
	RoleName string `json:"role_name" example:"Administrador"` // Nombre del rol del usuario
}
