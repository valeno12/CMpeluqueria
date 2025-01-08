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

// @Summary Crear turno
// @Description Permite crear un nuevo turno en el sistema.
// @Tags Turnos
// @Accept json
// @Produce json
// @Param request body dtos.CreateAppointmentDto true "Datos del turno"
// @Success 200 {object} dtos.Response{message=string,data=nil} "Turno creado con éxito"
// @Failure 400 {object} dtos.Response{message=string,data=nil} "Datos inválidos"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno [post]
// @Security BearerAuth
func CreateAppointment(c echo.Context) error {
	var appointment dtos.CreateAppointmentDto
	if err := c.Bind(&appointment); err != nil {
		logger.Log.Warn("[AppointmentController][CreateAppointment] Error: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Los datos enviados son inválidos: "+err.Error())
	}

	if err := services.CreateAppointment(appointment); err != nil {
		logger.Log.Error("[AppointmentController][CreateAppointment] Error al crear turno: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo crear el turno: "+err.Error())
	}

	logger.Log.Info("[AppointmentController][CreateAppointment] Turno creado exitosamente")
	return helpers.RespondSuccess(c, "Turno creado con éxito", nil)
}

// @Summary Obtener todos los turnos
// @Description Devuelve una lista de todos los turnos registrados, con la posibilidad de aplicar filtros.
// @Tags Turnos
// @Produce json
// @Param client_id query string false "ID del cliente (opcional)"
// @Param status query string false "Estado del turno (opcional)"
// @Param start_date query string false "Fecha de inicio (opcional), formato: YYYY-MM-DD"
// @Param end_date query string false "Fecha de fin (opcional), formato: YYYY-MM-DD"
// @Success 200 {object} dtos.Response{message=string,data=[]dtos.AllAppointmentDto} "Turnos obtenidos"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno [get]
// @Security BearerAuth
func GetAllAppointments(c echo.Context) error {
	clientID := c.QueryParam("client_id")
	status := c.QueryParam("status")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")

	appointments, err := services.GetAllAppointments(clientID, status, startDate, endDate)
	if err != nil {
		logger.Log.Error("[AppointmentController][GetAllAppointments] Error al obtener turnos: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.RespondSuccess(c, "Turnos obtenidos", appointments)
}

// @Summary Obtener turno por ID
// @Description Devuelve los datos de un turno específico.
// @Tags Turnos
// @Produce json
// @Param id path int true "ID del turno"
// @Success 200 {object} dtos.Response{message=string,data=dtos.AppointmentByIDDto} "Turno obtenido con éxito"
// @Failure 400 {object} dtos.Response{message=string,data=nil} "ID inválido"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno/{id} [get]
// @Security BearerAuth
func GetAppointmentByID(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][GetAppointmentByID] Error: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "El ID del turno es inválido")
	}

	appointment, err := services.GetAppointmentByID(uint(appointmentID))
	if err != nil {
		logger.Log.Error("[AppointmentController][GetAppointmentByID] Error al obtener turno con ID: ", appointmentID, " - ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo obtener el turno")
	}

	logger.Log.Infof("[AppointmentController][GetAppointmentByID] Turno obtenido con ID: %d", appointmentID)
	return helpers.RespondSuccess(c, "Turno obtenido con éxito", appointment)
}

// @Summary Actualizar turno
// @Description Permite actualizar un turno existente en el sistema.
// @Tags Turnos
// @Accept json
// @Produce json
// @Param id path int true "ID del turno"
// @Param request body dtos.CreateAppointmentDto true "Datos actualizados del turno"
// @Success 200 {object} dtos.Response{message=string,data=nil} "Turno actualizado con éxito"
// @Failure 400 {object} dtos.Response{message=string,data=nil} "Datos o ID inválidos"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno/{id} [put]
// @Security BearerAuth
func UpdateAppointment(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointment] Error: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "El ID del turno es inválido")
	}

	var appointmentDto dtos.CreateAppointmentDto
	if err := c.Bind(&appointmentDto); err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointment] Error: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Los datos enviados son inválidos")
	}

	if err := services.UpdateAppointment(uint(appointmentID), appointmentDto); err != nil {
		logger.Log.Error("[AppointmentController][UpdateAppointment] Error al actualizar turno con ID: ", appointmentID, " - ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo actualizar el turno")
	}

	logger.Log.Infof("[AppointmentController][UpdateAppointment] Turno actualizado con ID: %d", appointmentID)
	return helpers.RespondSuccess(c, "Turno actualizado con éxito", nil)
}

// @Summary Actualizar productos del turno
// @Description Permite actualizar los productos utilizados en un turno finalizado.
// @Tags Turnos
// @Accept json
// @Produce json
// @Param id path int true "ID del turno"
// @Param request body dtos.UpdateAppointmentProductsDto true "Datos de productos"
// @Success 200 {object} dtos.Response{message=string,data=nil} "Productos actualizados"
// @Failure 400 {object} dtos.Response{message=string,data=nil} "Datos inválidos"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno/{id}/products [put]
// @Security BearerAuth
func UpdateAppointmentProducts(c echo.Context) error {
	idParam := c.Param("id")
	appointmentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointmentProducts] Error: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "El ID del turno es inválido")
	}

	var productsDto dtos.UpdateAppointmentProductsDto
	if err := c.Bind(&productsDto); err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointmentProducts] Error: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Los datos enviados son inválidos")
	}

	if err := services.UpdateAppointmentProducts(uint(appointmentID), productsDto); err != nil {
		logger.Log.Error("[AppointmentController][UpdateAppointmentProducts] Error al actualizar productos del turno con ID: ", appointmentID, " - ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo actualizar los productos del turno")
	}

	logger.Log.Infof("[AppointmentController][UpdateAppointmentProducts] Productos actualizados para turno ID: %d", appointmentID)
	return helpers.RespondSuccess(c, "Productos del turno actualizados con éxito", nil)
}

// @Summary Finalizar turno
// @Description Permite finalizar un turno, registrando productos utilizados y método de pago.
// @Tags Turnos
// @Accept json
// @Produce json
// @Param id path int true "ID del turno"
// @Param request body dtos.FinalizeAppointmentDto true "Datos para finalizar el turno"
// @Success 200 {object} dtos.Response{message=string,data=nil} "Turno finalizado con éxito"
// @Failure 400 {object} dtos.Response{message=string,data=nil} "Datos o ID inválidos"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno/{id}/finalizar [put]
// @Security BearerAuth
func FinalizeAppointment(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][FinalizeAppointment] Error: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "El ID del turno es inválido")
	}

	var finalizeDto dtos.FinalizeAppointmentDto
	if err := c.Bind(&finalizeDto); err != nil {
		logger.Log.Warn("[AppointmentController][FinalizeAppointment] Error: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Los datos enviados son inválidos")
	}

	if err := services.FinalizeAppointment(uint(appointmentID), finalizeDto); err != nil {
		logger.Log.Error("[AppointmentController][FinalizeAppointment] Error al finalizar turno con ID: ", appointmentID, " - ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo finalizar el turno: "+err.Error())
	}

	logger.Log.Infof("[AppointmentController][FinalizeAppointment] Turno finalizado con éxito: ID %d", appointmentID)
	return helpers.RespondSuccess(c, "Turno finalizado con éxito", nil)
}

// @Summary Eliminar turno
// @Description Permite eliminar un turno del sistema.
// @Tags Turnos
// @Produce json
// @Param id path int true "ID del turno"
// @Success 200 {object} dtos.Response{message=string,data=nil} "Turno eliminado con éxito"
// @Failure 400 {object} dtos.Response{message=string,data=nil} "ID inválido"
// @Failure 500 {object} dtos.Response{message=string,data=nil} "Error interno del servidor"
// @Router /turno/{id} [delete]
// @Security BearerAuth
func DeleteAppointment(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][DeleteAppointment] Error: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "El ID del turno es inválido")
	}

	if err := services.DeleteAppointment(uint(appointmentID)); err != nil {
		logger.Log.Error("[AppointmentController][DeleteAppointment] Error al eliminar turno con ID: ", appointmentID, " - ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo eliminar el turno")
	}

	logger.Log.Infof("[AppointmentController][DeleteAppointment] Turno eliminado con éxito: ID %d", appointmentID)
	return helpers.RespondSuccess(c, "Turno eliminado con éxito", nil)
}
