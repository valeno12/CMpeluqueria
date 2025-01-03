package database

import (
	"peluqueria/internal/models"
	"peluqueria/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase se asegura de que la base de datos tenga datos iniciales.
func SeedDatabase(db *gorm.DB) error {
	logger.Log.Info("Iniciando el proceso de seeders...")
	seedRoles(db)
	seedUsers(db)
	seedPermissions(db)
	seedRolePermissions(db)
	logger.Log.Info("Seeders ejecutados con éxito")
	return nil
}

func seedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "empleado"},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Solo inserta el rol si no existe
				if err := db.Create(&role).Error; err != nil {
					logger.Log.Error("Error al insertar rol '", role.Name, "': ", err)
				} else {
					logger.Log.Info("Rol '", role.Name, "' creado con éxito")
				}
			}
		}
	}
}

func seedUsers(db *gorm.DB) {
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		logger.Log.Error("Error al buscar rol admin: ", err)
		return
	}

	adminUser := models.User{
		Username: "admin",
		Password: HashPassword("admin123"),
		RoleID:   adminRole.ID,
	}

	var existingUser models.User
	if err := db.Where("username = ?", adminUser.Username).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&adminUser).Error; err != nil {
				logger.Log.Error("Error al crear usuario admin: ", err)
			} else {
				logger.Log.Info("Usuario admin creado con éxito")
			}
		}
	}
}

func seedPermissions(db *gorm.DB) {
	permissions := []models.Permission{
		{Name: "create_appointment", Description: "Crear turnos"},
		{Name: "update_appointment", Description: "Editar turnos"},
		{Name: "delete_appointment", Description: "Eliminar turnos"},
		{Name: "create_service", Description: "Crear servicios"},
		{Name: "update_service", Description: "Editar servicios"},
		{Name: "delete_service", Description: "Eliminar servicios"},
		{Name: "create_product", Description: "Crear productos"},
		{Name: "update_product", Description: "Editar productos"},
		{Name: "delete_product", Description: "Eliminar productos"},
		{Name: "create_user", Description: "Crear usuarios"},
		{Name: "update_user", Description: "Editar usuarios"},
		{Name: "delete_user", Description: "Eliminar usuarios"},
		{Name: "create_role", Description: "Crear roles"},
		{Name: "update_role", Description: "Editar roles"},
		{Name: "delete_role", Description: "Eliminar roles"},
		{Name: "create_client", Description: "Crear clientes"},
		{Name: "update_client", Description: "Editar clientes"},
		{Name: "delete_client", Description: "Eliminar clientes"},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := db.Where("name = ?", permission.Name).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					logger.Log.Error("Error al insertar permiso '", permission.Name, "': ", err)
				} else {
					logger.Log.Info("Permiso '", permission.Name, "' creado con éxito")
				}
			}
		}
	}
}

func seedRolePermissions(db *gorm.DB) {
	rolePermissions := map[string][]string{
		"admin": {
			"create_appointment", "update_appointment", "delete_appointment",
			"create_service", "update_service", "delete_service",
			"create_product", "update_product", "delete_product",
			"create_user", "update_user", "delete_user",
			"create_role", "update_role", "delete_role", "create_client", "update_client", "delete_client",
		},
		"empleado": {
			"create_appointment", "update_appointment",
			"create_service", "update_service",
		},
	}

	for roleName, permissions := range rolePermissions {
		var role models.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			logger.Log.Error("Error al buscar rol '", roleName, "': ", err)
			continue
		}

		for _, permissionName := range permissions {
			var permission models.Permission
			if err := db.Where("name = ?", permissionName).First(&permission).Error; err != nil {
				logger.Log.Error("Error al buscar permiso '", permissionName, "': ", err)
				continue
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := db.FirstOrCreate(&rolePermission, rolePermission).Error; err != nil {
				logger.Log.Error("Error al asignar permiso '", permissionName, "' al rol '", roleName, "': ", err)
			} else {
				logger.Log.Info("Permiso '", permissionName, "' asignado al rol '", roleName, "' con éxito")
			}
		}
	}
}

// HashPassword es una función auxiliar para encriptar contraseñas
func HashPassword(pass string) string {
	costo := 8
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), costo)
	return string(bytes)
}
