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

// @Summary Crear rol
// @Description Crea un nuevo rol con permisos específicos.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtos.CreateRoleDto true "Datos del nuevo rol"
// @Success 200 {object} dtos.Response{data=nil} "Rol creado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Datos inválidos. Ejemplo: El nombre del rol ya existe"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /rol [post]
func CreateRole(c echo.Context) error {
	logger.Log.Info("[RoleController][CreateRole] Intentando crear un rol")
	var role dtos.CreateRoleDto
	if err := c.Bind(&role); err != nil {
		logger.Log.Warn("[RoleController][CreateRole] Error al crear rol: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err := services.CreateRole(role)
	if err != nil {
		logger.Log.Error("[RoleController][CreateRole] Error al crear rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[RoleController][CreateRole] Rol creado con éxito: %s", role.Name)
	return helpers.RespondSuccess(c, "Rol creado exitosamente", nil)
}

// @Summary Obtener todos los roles
// @Description Devuelve una lista de todos los roles registrados en el sistema.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.Response{data=[]dtos.GetRoleDto} "Lista de roles encontrados"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /rol [get]
func GetAllRoles(c echo.Context) error {
	logger.Log.Info("[RoleController][GetAllRoles] Intentando obtener todos los roles")
	roles, err := services.GetAllRoles()
	if err != nil {
		logger.Log.Error("[RoleController][GetAllRoles] Error al obtener roles: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[RoleController][GetAllRoles] Roles obtenidos: %d roles encontrados", len(roles))
	return helpers.RespondSuccess(c, "Roles obtenidos", roles)
}

// @Summary Obtener rol por ID
// @Description Devuelve los datos de un rol específico.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del rol"
// @Success 200 {object} dtos.Response{data=dtos.GetRoleDto} "Rol encontrado"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido"
// @Failure 404 {object} dtos.ErrorResponse "Rol no encontrado"
// @Router /rol/{id} [get]
func GetRoleByID(c echo.Context) error {
	id := c.Param("id")
	logger.Log.Infof("[RoleController][GetRoleByID] Intentando obtener rol con ID: %s", id)
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[RoleController][GetRoleByID] Error al obtener rol: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	role, err := services.GetRoleByID(uint(roleID))
	if err != nil {
		logger.Log.Error("[RoleController][GetRoleByID] Error al obtener rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[RoleController][GetRoleByID] Rol encontrado: %s", role.Name)
	return helpers.RespondSuccess(c, "Rol encontrado", role)
}

// @Summary Actualizar rol
// @Description Actualiza los datos de un rol específico.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del rol"
// @Param request body dtos.CreateRoleDto true "Datos actualizados del rol"
// @Success 200 {object} dtos.Response{data=nil} "Rol actualizado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Datos inválidos o ID inválido"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /rol/{id} [put]
func UpdateRole(c echo.Context) error {
	id := c.Param("id")
	logger.Log.Infof("[RoleController][UpdateRole] Intentando actualizar rol con ID: %s", id)
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[RoleController][UpdateRole] Error al actualizar rol: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	var role dtos.CreateRoleDto
	if err := c.Bind(&role); err != nil {
		logger.Log.Warn("[RoleController][UpdateRole] Error al actualizar rol: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err = services.UpdateRole(uint(roleID), role)
	if err != nil {
		logger.Log.Error("[RoleController][UpdateRole] Error al actualizar rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[RoleController][UpdateRole] Rol actualizado con éxito: %s", role.Name)
	return helpers.RespondSuccess(c, "Rol actualizado correctamente", nil)
}

// @Summary Eliminar rol
// @Description Elimina un rol específico.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del rol"
// @Success 200 {object} dtos.Response{data=nil} "Rol eliminado correctamente"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /rol/{id} [delete]
func DeleteRole(c echo.Context) error {
	id := c.Param("id")
	logger.Log.Infof("[RoleController][DeleteRole] Intentando eliminar rol con ID: %s", id)
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[RoleController][DeleteRole] Error al eliminar rol: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	err = services.DeleteRole(uint(roleID))
	if err != nil {
		logger.Log.Error("[RoleController][DeleteRole] Error al eliminar rol: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[RoleController][DeleteRole] Rol eliminado con éxito: ID %d", roleID)
	return helpers.RespondSuccess(c, "Rol eliminado correctamente", nil)
}
