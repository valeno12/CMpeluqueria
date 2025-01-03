package controllers

import (
	"net/http"
	"peluqueria/internal/dtos"
	"peluqueria/internal/helpers"
	"peluqueria/internal/services"
	"peluqueria/logger"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateRole(c echo.Context) error {
	var role dtos.CreateRoleDto
	if err := c.Bind(&role); err != nil {
		logger.Log.Warn("Error al crear usuario: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err := services.CreateRole(role)
	if err != nil {
		logger.Log.Error("Error al crear usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("Rol creado: %s", role.Name)
	return helpers.RespondSuccess(c, "Rol creado", nil)
}

func GetAllRoles(c echo.Context) error {
	roles, err := services.GetAllRoles()
	if err != nil {
		logger.Log.Error("Error al obtener roles: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.RespondSuccess(c, "Roles obtenidos", roles)
}

func GetRoleByID(c echo.Context) error {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("Error al obtener rol: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	user, err := services.GetRoleByID(uint(roleID))
	if err != nil {
		logger.Log.Error("Error al obtener rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Info("Rol encontrado: ", user.Name)
	return helpers.RespondSuccess(c, "Rol encontrado", user)
}

func UpdateRole(c echo.Context) error {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("Error al actualizar rol: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	var role dtos.CreateRoleDto
	if err := c.Bind(&role); err != nil {
		logger.Log.Warn("Error al actualizar rol: datos inválidos " + err.Error())
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err = services.UpdateRole(uint(roleID), role)
	if err != nil {
		logger.Log.Error("Error al actualizar rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("Rol actualizado: %s", role.Name)
	return helpers.RespondSuccess(c, "Rol actualizado", nil)
}

func DeleteRole(c echo.Context) error {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("Error al eliminar rol: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	err = services.DeleteRole(uint(roleID))
	if err != nil {
		logger.Log.Error("Error al eliminar rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Info("Rol eliminado: ID ", roleID)
	return helpers.RespondSuccess(c, "Rol eliminado", nil)
}
