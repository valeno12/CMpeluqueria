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

// CreateService crea un nuevo servicio.
// @Summary Crear servicio
// @Description Permite crear un nuevo servicio en el sistema.
// @Tags Servicios
// @Accept json
// @Produce json
// @Param request body dtos.ServiceDto true "Datos del servicio"
// @Success 200 {object} dtos.Response{data=nil} "Servicio creado"
// @Failure 400 {object} dtos.ErrorResponse "Datos inválidos"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /servicio [post]
// @Security BearerAuth
func CreateService(c echo.Context) error {
	logger.Log.Info("[ServiceController][CreateService] Iniciando creación de servicio")

	var service dtos.ServiceDto
	if err := c.Bind(&service); err != nil {
		logger.Log.Warn("[ServiceController][CreateService] Error al parsear datos: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err := services.CreateService(service)
	if err != nil {
		logger.Log.Error("[ServiceController][CreateService] Error al crear servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[ServiceController][CreateService] Servicio creado con éxito: %s", service.Name)
	return helpers.RespondSuccess(c, "Servicio creado", nil)
}

// GetAllServices obtiene la lista de servicios.
// @Summary Obtener todos los servicios
// @Description Devuelve una lista de todos los servicios registrados en el sistema.
// @Tags Servicios
// @Produce json
// @Success 200 {object} dtos.Response{data=[]dtos.GetServiceDto} "Servicios obtenidos"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /servicio [get]
// @Security BearerAuth
func GetAllServices(c echo.Context) error {
	logger.Log.Info("[ServiceController][GetAllServices] Obteniendo lista de servicios")

	services, err := services.GetAllServices()
	if err != nil {
		logger.Log.Error("[ServiceController][GetAllServices] Error al obtener servicios: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[ServiceController][GetAllServices] %d servicios obtenidos", len(services))
	return helpers.RespondSuccess(c, "Servicios obtenidos", services)
}

// GetServiceByID obtiene un servicio por su ID.
// @Summary Obtener servicio por ID
// @Description Devuelve los datos de un servicio específico.
// @Tags Servicios
// @Produce json
// @Param id path int true "ID del servicio" example:"10"
// @Success 200 {object} dtos.Response{data=dtos.GetServiceDto} "Servicio obtenido"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /servicio/{id} [get]
// @Security BearerAuth
func GetServiceByID(c echo.Context) error {
	logger.Log.Info("[ServiceController][GetServiceByID] Iniciando búsqueda de servicio por ID")

	id := c.Param("id")
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ServiceController][GetServiceByID] ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	service, err := services.GetServiceByID(uint(serviceID))
	if err != nil {
		logger.Log.Error("[ServiceController][GetServiceByID] Error al obtener servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[ServiceController][GetServiceByID] Servicio obtenido con éxito: %s", service.Name)
	return helpers.RespondSuccess(c, "Servicio obtenido", service)
}

// UpdateService actualiza un servicio.
// @Summary Actualizar servicio
// @Description Permite actualizar los datos de un servicio existente.
// @Tags Servicios
// @Accept json
// @Produce json
// @Param id path int true "ID del servicio" example:"10"
// @Param request body dtos.ServiceDto true "Datos del servicio"
// @Success 200 {object} dtos.Response{data=nil} "Servicio actualizado"
// @Failure 400 {object} dtos.ErrorResponse "ID o datos inválidos"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /servicio/{id} [put]
// @Security BearerAuth
func UpdateService(c echo.Context) error {
	logger.Log.Info("[ServiceController][UpdateService] Iniciando actualización de servicio")

	id := c.Param("id")
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ServiceController][UpdateService] ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var service dtos.ServiceDto
	if err := c.Bind(&service); err != nil {
		logger.Log.Warn("[ServiceController][UpdateService] Error al parsear datos: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.UpdateService(uint(serviceID), service)
	if err != nil {
		logger.Log.Error("[ServiceController][UpdateService] Error al actualizar servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[ServiceController][UpdateService] Servicio actualizado con éxito: %s", service.Name)
	return helpers.RespondSuccess(c, "Servicio actualizado", nil)
}

// DeleteService elimina un servicio.
// @Summary Eliminar servicio
// @Description Permite eliminar un servicio específico.
// @Tags Servicios
// @Produce json
// @Param id path int true "ID del servicio" example:"10"
// @Success 200 {object} dtos.Response{data=nil} "Servicio eliminado"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /servicio/{id} [delete]
// @Security BearerAuth
func DeleteService(c echo.Context) error {
	logger.Log.Info("[ServiceController][DeleteService] Iniciando eliminación de servicio")

	id := c.Param("id")
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ServiceController][DeleteService] ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	err = services.DeleteService(uint(serviceID))
	if err != nil {
		logger.Log.Error("[ServiceController][DeleteService] Error al eliminar servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Info("[ServiceController][DeleteService] Servicio eliminado con éxito")
	return helpers.RespondSuccess(c, "Servicio eliminado", nil)
}
