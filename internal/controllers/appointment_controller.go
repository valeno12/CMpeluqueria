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

func CreateAppointment(c echo.Context) error {
	var Appointment dtos.CreateAppointmentDto
	if err := c.Bind(&Appointment); err != nil {
		logger.Log.Warn("[AppointmentController][CreateAppointment] Error al crear cita: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "datos invalidos "+err.Error())
	}
	err := services.CreateAppointment(Appointment)
	if err != nil {
		logger.Log.Error("[AppointmentController][CreateAppointment] Error al crear cita: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[AppointmentController][CreateAppointment] Cita creada con exito")
	return helpers.RespondSuccess(c, "Cita creada", nil)
}

func GetAllAppointments(c echo.Context) error {
	appointments, err := services.GetAllAppointments()
	if err != nil {
		logger.Log.Error("[AppointmentController][GetAllAppointments] Error al obtener turnos: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.RespondSuccess(c, "Turnos obtenidas", appointments)
}

func GetAppointmentByID(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][GetAppointmentById] Error al obtener turno: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	appointment, err := services.GetAppointmentByID(uint(appointmentID))
	if err != nil {
		logger.Log.Error("[AppointmentController][GetAppointmentById] Error al obtener turno: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.RespondSuccess(c, "Turno obtenida", appointment)
}

func UpdateAppointment(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointment] Error al obtener turno: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	var appointmentDto dtos.CreateAppointmentDto
	if err := c.Bind(&appointmentDto); err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointment] Error: datos invalidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos invalidos")
	}
	err = services.UpdateAppointment(uint(appointmentID), appointmentDto)
	if err != nil {
		logger.Log.Error("[AppointmentController][UpdateAppointment] Error al actualizar turno")
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.RespondSuccess(c, "Turno actualizado", nil)
}

func DeleteAppointment(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][DeleteAppointment] Error al obtener turno: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	err = services.DeleteAppointment(uint(appointmentID))
	if err != nil {
		logger.Log.Error("[AppointmentController][DeleteAppointment] Error al eliminar turno: " + err.Error())
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("Turno eliminado correctamente")
	return helpers.RespondSuccess(c, "Turno eliminado correctamente", nil)
}

func FinalizeAppointment(c echo.Context) error {
	id := c.Param("id")
	appointmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][FinalizeAppointment] Error al finalizar turno: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var finalizeDto dtos.FinalizeAppointmentDto
	if err := c.Bind(&finalizeDto); err != nil {
		logger.Log.Warn("[AppointmentController][FinalizeAppointment] Error al finalizar turno: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.FinalizeAppointment(uint(appointmentID), finalizeDto)
	if err != nil {
		logger.Log.Error("[AppointmentController][FinalizeAppointment] Error al finalizar turno: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[AppointmentController][FinalizeAppointment] Turno finalizado con éxito: ID %d", appointmentID)
	return helpers.RespondSuccess(c, "Turno finalizado con éxito", nil)
}

func UpdateAppointmentProducts(c echo.Context) error {
	idParam := c.Param("id")
	appointmentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointmentProducts] ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var productsDto dtos.UpdateAppointmentProductsDto
	if err := c.Bind(&productsDto); err != nil {
		logger.Log.Warn("[AppointmentController][UpdateAppointmentProducts] Datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.UpdateAppointmentProducts(uint(appointmentID), productsDto)
	if err != nil {
		logger.Log.Error("[AppointmentController][UpdateAppointmentProducts] Error al actualizar productos: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[AppointmentController][UpdateAppointmentProducts] Productos actualizados para turno ID: %d", appointmentID)
	return helpers.RespondSuccess(c, "Productos actualizados con éxito", nil)
}
