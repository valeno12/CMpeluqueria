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

func CreateService(c echo.Context) error {

	var service dtos.ServiceDto
	if err := c.Bind(&service); err != nil {
		logger.Log.Warn("[ServiceController][CreateService] Error al crear servicio: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err := services.CreateService(service)
	if err != nil {
		logger.Log.Error("[ServiceController][CreateService] Error al crear servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ServiceController][CreateService] Servicio creado: %s", service.Name)
	return helpers.RespondSuccess(c, "Servicio creado", nil)
}

func GetAllServices(c echo.Context) error {
	services, err := services.GetAllServices()
	if err != nil {
		logger.Log.Error("[ServiceController][GetAllServices] Error al obtener servicios: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Info("[ServiceController][GetAllServices] Servicios obtenidos")
	return helpers.RespondSuccess(c, "Servicios obtenidos", services)
}

func GetServiceByID(c echo.Context) error {
	id := c.Param("id")
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ServiceController][GetServiceByID] Error al actualizar producto: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	service, err := services.GetServiceByID(uint(serviceID))
	if err != nil {
		logger.Log.Error("[ServiceController][GetServiceByID] Error al obtener servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ServiceController][GetServiceByID] Servicio obtenido: %s", service.Name)
	return helpers.RespondSuccess(c, "Servicio obtenido", service)
}

func UpdateService(c echo.Context) error {
	id := c.Param("id")
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ServiceController][UpdateService] Error al actualizar servicio: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var service dtos.ServiceDto
	if err := c.Bind(&service); err != nil {
		logger.Log.Warn("[ServiceController][UpdateService] Error al actualizar servicio: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.UpdateService(uint(serviceID), service)
	if err != nil {
		logger.Log.Error("[ServiceController][UpdateService] Error al actualizar servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ServiceController][UpdateService] Servicio actualizado: %s", service.Name)
	return helpers.RespondSuccess(c, "Servicio actualizado", nil)
}

func DeleteService(c echo.Context) error {
	id := c.Param("id")
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ServiceController][DeleteService] Error al eliminar servicio: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	err = services.DeleteService(uint(serviceID))
	if err != nil {
		logger.Log.Error("[ServiceController][DeleteService] Error al eliminar servicio: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Info("[ServiceController][DeleteService] Servicio eliminado")
	return helpers.RespondSuccess(c, "Servicio eliminado", nil)
}
