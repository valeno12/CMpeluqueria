package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"

	"gorm.io/gorm"
)

func CreateRole(roleDto dtos.CreateRoleDto) error {
	logger.Log.Infof("[RoleService][CreateRole] Intentando crear rol: %s", roleDto.Name)

	// Verificar si el rol ya existe
	var existingRole models.Role
	if err := database.DB.Where("name = ?", roleDto.Name).First(&existingRole).Error; err == nil {
		logger.Log.Warnf("[RoleService][CreateRole] El rol ya existe: %s", roleDto.Name)
		return errors.New("el nombre del rol ya existe")
	}

	role := models.Role{
		Name: roleDto.Name,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		logger.Log.Error("[RoleService][CreateRole] Error al crear el rol: ", err)
		return errors.New("error al crear el rol")
	}
	logger.Log.Infof("[RoleService][CreateRole] Rol creado con éxito: %s (ID: %d)", role.Name, role.ID)

	// Asignar permisos si hay nombres en PermissionNames
	if len(roleDto.PermissionNames) > 0 {
		logger.Log.Infof("[RoleService][CreateRole] Asignando permisos al rol: %s", role.Name)
		for _, permissionName := range roleDto.PermissionNames {
			// Validar si el permiso existe
			var permission models.Permission
			if err := database.DB.Where("name = ?", permissionName).First(&permission).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.Log.Warnf("[RoleService][CreateRole] Permiso no encontrado: %s", permissionName)
					return errors.New("uno o más permisos no existen")
				}
				logger.Log.Error("[RoleService][CreateRole] Error al buscar permiso: ", err)
				return err
			}

			// Crear la relación en RolePermission
			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := database.DB.Create(&rolePermission).Error; err != nil {
				logger.Log.Error("[RoleService][CreateRole] Error al asignar permisos al rol: ", err)
				return errors.New("error al asignar permisos al rol")
			}
			logger.Log.Infof("[RoleService][CreateRole] Permiso asignado: RoleID %d, PermissionID %d", role.ID, permission.ID)
		}
	}
	logger.Log.Infof("[RoleService][CreateRole] Permisos asignados con éxito al rol: %s", role.Name)
	return nil
}

func GetAllRoles() ([]dtos.GetRoleDto, error) {
	logger.Log.Info("[RoleService][GetAllRoles] Intentando obtener roles")

	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		logger.Log.Error("[RoleService][GetAllRoles] Error al obtener roles: ", err)
		return nil, errors.New("error al obtener roles")
	}

	rolesDto := make([]dtos.GetRoleDto, len(roles))
	for i, role := range roles {
		permissions := make([]string, len(role.Permissions))
		for j, permission := range role.Permissions {
			permissions[j] = permission.Name
		}
		rolesDto[i] = dtos.GetRoleDto{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		}
	}

	logger.Log.Infof("[RoleService][GetAllRoles] Roles obtenidos: %d", len(roles))
	return rolesDto, nil
}

func GetRoleByID(id uint) (dtos.GetRoleDto, error) {
	logger.Log.Infof("[RoleService][GetRoleByID] Intentando obtener rol: ID %d", id)

	var role models.Role

	err := database.DB.Preload("Permissions").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[RoleService][GetRoleByID] Rol no encontrado: ID %d", id)
			return dtos.GetRoleDto{}, errors.New("rol no encontrado")
		}
		logger.Log.Error("[RoleService][GetRoleByID] Error al obtener rol: ", err)
		return dtos.GetRoleDto{}, errors.New("error al obtener rol")
	}

	permissions := make([]string, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissions[i] = permission.Name
	}

	roleDto := dtos.GetRoleDto{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}

	logger.Log.Infof("[RoleService][GetRoleByID] Rol obtenido con éxito: %s (ID: %d)", role.Name, role.ID)
	return roleDto, nil
}

func UpdateRole(roleID uint, roleDto dtos.CreateRoleDto) error {
	logger.Log.Infof("[RoleService][UpdateRole] Actualizando rol con ID: %d", roleID)

	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[RoleService][UpdateRole] Rol no encontrado: ID %d", roleID)
			return errors.New("el rol no existe")
		}
		logger.Log.Error("[RoleService][UpdateRole] Error al buscar rol: ", err)
		return err
	}

	if roleDto.Name != "" {
		role.Name = roleDto.Name
	}

	if err := database.DB.Save(&role).Error; err != nil {
		logger.Log.Error("[RoleService][UpdateRole] Error al actualizar rol: ", err)
		return errors.New("error al actualizar el rol")
	}

	logger.Log.Infof("[RoleService][UpdateRole] Rol actualizado: %s (ID: %d)", role.Name, role.ID)

	// Manejar permisos
	if len(roleDto.PermissionNames) > 0 {
		logger.Log.Infof("[RoleService][UpdateRole] Actualizando permisos del rol: %s", role.Name)
		if err := database.DB.Unscoped().Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			logger.Log.Error("[RoleService][UpdateRole] Error al eliminar permisos existentes del rol: ", err)
			return errors.New("error al actualizar los permisos del rol")
		}

		for _, permissionName := range roleDto.PermissionNames {
			var permission models.Permission
			if err := database.DB.Where("name = ?", permissionName).First(&permission).Error; err != nil {
				logger.Log.Warnf("[RoleService][UpdateRole] Permiso no encontrado: %s", permissionName)
				return errors.New("uno o más permisos no existen")
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := database.DB.Create(&rolePermission).Error; err != nil {
				logger.Log.Error("[RoleService][UpdateRole] Error al asignar permisos al rol: ", err)
				return errors.New("error al actualizar los permisos del rol")
			}
			logger.Log.Infof("[RoleService][UpdateRole] Permiso asignado: RoleID %d, PermissionID %d", role.ID, permission.ID)
		}
		logger.Log.Infof("[RoleService][UpdateRole] Permisos actualizados para el rol: %s", role.Name)
	}

	return nil
}

func DeleteRole(id uint) error {
	logger.Log.Infof("[RoleService][DeleteRole] Intentando eliminar rol: ID %d", id)

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[RoleService][DeleteRole] Rol no encontrado: ID %d", id)
			return errors.New("rol no encontrado")
		}
		logger.Log.Error("[RoleService][DeleteRole] Error al buscar rol: ", err)
		return errors.New("error al buscar rol")
	}

	var usersCount int64
	if err := database.DB.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&usersCount).Error; err != nil {
		logger.Log.Error("[RoleService][DeleteRole] Error al verificar usuarios asignados al rol: ", err)
		return errors.New("error al verificar usuarios asignados al rol")
	}

	if usersCount > 0 {
		logger.Log.Warnf("[RoleService][DeleteRole] No se puede eliminar el rol porque está asignado a usuarios: %d usuarios", usersCount)
		return errors.New("no se puede eliminar el rol porque está asignado a usuarios")
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		logger.Log.Error("[RoleService][DeleteRole] Error al eliminar rol: ", err)
		return errors.New("error al eliminar rol")
	}

	if err := database.DB.Unscoped().Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
		logger.Log.Error("[RoleService][DeleteRole] Error al eliminar permisos asociados al rol: ", err)
		return errors.New("error al eliminar permisos asociados al rol")
	}

	logger.Log.Infof("[RoleService][DeleteRole] Rol eliminado con éxito: ID %d", id)
	return nil
}
