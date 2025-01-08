package dtos

type RoleDto struct {
	Name string `json:"name"`
}

type CreateRoleDto struct {
	Name            string   `json:"name" example:"empleado"`
	PermissionNames []string `json:"permission_names" example:"create_user, delete_user"`
}

type GetRoleDto struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name" example:"empleado"`
	Permissions []string `json:"permissions" example:"create_user, delete_user"` // Solo los nombres de los permisos
}
