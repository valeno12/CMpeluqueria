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
	logger.Log.Infof("Intentando crear rol: %s", roleDto.Name)

	// Verificar si el rol ya existe
	var existingRole models.Role
	if err := database.DB.Where("name = ?", roleDto.Name).First(&existingRole).Error; err == nil {
		logger.Log.Warnf("El rol ya existe: %s", roleDto.Name)
		return errors.New("el nombre del rol ya existe")
	}

	role := models.Role{
		Name: roleDto.Name,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		logger.Log.Error("Error al crear el rol: ", err)
		return errors.New("error al crear el rol")
	}
	logger.Log.Infof("Rol creado con éxito: %s (ID: %d)", role.Name, role.ID)

	// Asignar permisos si hay nombres en PermissionNames
	if len(roleDto.PermissionNames) > 0 {
		logger.Log.Infof("Asignando permisos al rol: %s", role.Name)
		for _, permissionName := range roleDto.PermissionNames {
			// Validar si el permiso existe
			var permission models.Permission
			if err := database.DB.Where("name = ?", permissionName).First(&permission).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.Log.Warnf("Permiso no encontrado: %s", permissionName)
					return errors.New("uno o más permisos no existen")
				}
				logger.Log.Error("Error al buscar permiso: ", err)
				return err
			}

			// Crear la relación en RolePermission
			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := database.DB.Create(&rolePermission).Error; err != nil {
				logger.Log.Error("Error al asignar permisos al rol: ", err)
				return errors.New("error al asignar permisos al rol")
			}
			logger.Log.Infof("Permiso asignado: RoleID %d, PermissionID %d", role.ID, permission.ID)
		}
	}
	logger.Log.Infof("Permisos asignados con éxito al rol: %s", role.Name)
	return nil
}

func GetAllRoles() ([]models.Role, error) {
	logger.Log.Info("Intentando obtener roles")

	var roles []models.Role
	if err := database.DB.Find(&roles).Error; err != nil {
		logger.Log.Error("Error al obtener roles: ", err)
		return nil, errors.New("error al obtener roles")
	}

	logger.Log.Infof("Roles obtenidos: %d", len(roles))
	return roles, nil
}

func GetRoleByID(id uint) (dtos.CreateRoleDto, error) {
	logger.Log.Infof("Intentando obtener rol: ID %d", id)

	var role models.Role

	err := database.DB.Preload("Permissions").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("Rol no encontrado: ID %d", id)
			return dtos.CreateRoleDto{}, errors.New("rol no encontrado")
		}
		logger.Log.Error("Error al obtener rol: ", err)
		return dtos.CreateRoleDto{}, errors.New("error al obtener rol")
	}

	permissionNames := make([]string, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissionNames[i] = permission.Name
	}

	return dtos.CreateRoleDto{
		Name:            role.Name,
		PermissionNames: permissionNames,
	}, nil
}

func UpdateRole(roleID uint, roleDto dtos.CreateRoleDto) error {
	logger.Log.Infof("Actualizando rol con ID: %d", roleID)
	if roleID == 0 {
		logger.Log.Warn("ID del rol faltante en actualización")
		return errors.New("el ID del rol es obligatorio")
	}

	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Log.Warnf("rol no encontrado para actualizar: ID %d", roleID)
			return errors.New("el rol no existe")
		}
		logger.Log.Error("Error al buscar rol para actualizar: ", err)
		return err
	}

	if roleDto.Name != "" {
		role.Name = roleDto.Name
	}

	if err := database.DB.Save(&role).Error; err != nil {
		logger.Log.Error("Error al actualizar el rol: ", err)
		return errors.New("error al actualizar el rol")
	}

	logger.Log.Infof("Nombre del rol actualizado: %s", role.Name)

	// Manejar permisos si se proporcionaron
	if len(roleDto.PermissionNames) > 0 {
		logger.Log.Infof("Actualizando permisos del rol: %s", role.Name)
		// Eliminar relaciones existentes en RolePermission
		if err := database.DB.Unscoped().Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			logger.Log.Error("Error al eliminar permisos existentes del rol: ", err)
			return errors.New("error al actualizar los permisos del rol " + err.Error())
		}
		// Crear nuevas relaciones de permisos
		for _, permissionName := range roleDto.PermissionNames {
			var permission models.Permission
			if err := database.DB.Where("name = ?", permissionName).First(&permission).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.Log.Warnf("Permiso no encontrado: %s", permissionName)
					return errors.New("uno o más permisos no existen")
				}
				logger.Log.Error("Error al buscar permiso: ", err)
				return err
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := database.DB.Create(&rolePermission).Error; err != nil {
				logger.Log.Error("Error al asignar permisos al rol: ", err)
				return errors.New("error al actualizar los permisos del rol")
			}
			logger.Log.Infof("Permiso asignado: RoleID %d, PermissionID %d", role.ID, permission.ID)
		}
		logger.Log.Infof("Permisos actualizados con éxito para el rol: %s", role.Name)
	}

	return nil
}

func DeleteRole(id uint) error {
	logger.Log.Infof("Intentando eliminar rol: ID %d", id)

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("Rol no encontrado: ID %d", id)
			return errors.New("rol no encontrado")
		}
		logger.Log.Error("Error al buscar rol: ", err)
		return errors.New("error al buscar rol")
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		logger.Log.Error("Error al eliminar rol: ", err)
		return errors.New("error al eliminar rol")
	}

	if err := database.DB.Unscoped().Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
		logger.Log.Error("Error al eliminar permisos existentes del rol: ", err)
		return errors.New("error al actualizar los permisos del rol " + err.Error())
	}
	return nil
}
