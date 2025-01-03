package dtos

type RoleDto struct {
	Name string `json:"name"`
}

type CreateRoleDto struct {
	Name            string   `json:"name"`
	PermissionNames []string `json:"permission_names"`
}
