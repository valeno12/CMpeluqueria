package dtos

type RoleDto struct {
	Name string `json:"name"`
}

type CreateRoleDto struct {
	Name            string   `json:"name"`
	PermissionNames []string `json:"permission_names"`
}

type GetRoleDto struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"` // Solo los nombres de los permisos
}
