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

// @Summary Crear cliente
// @Description Crea un nuevo cliente en el sistema.
// @Tags Clientes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtos.ClientDTO true "Datos del cliente"
// @Success 200 {object} dtos.Response{data=nil} "Cliente creado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Datos inválidos. Ejemplo: El cliente ya existe"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /cliente [post]
func CreateClient(c echo.Context) error {
	logger.Log.Info("[ClientController][CreateClient] Intentando crear cliente")
	var clientDTO dtos.ClientDTO
	if err := c.Bind(&clientDTO); err != nil {
		logger.Log.Warn("[ClientController][CreateClient] Error al crear cliente: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err := services.CreateClient(clientDTO)
	if err != nil {
		logger.Log.Error("[ClientController][CreateClient] Error al crear cliente: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ClientController][CreateClient] Cliente creado: %s", clientDTO.Name)
	return helpers.RespondSuccess(c, "Cliente creado exitosamente", nil)
}

// @Summary Obtener todos los clientes
// @Description Devuelve una lista de todos los clientes registrados.
// @Tags Clientes
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.Response{data=[]dtos.GetClientDto} "Lista de clientes encontrados"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /cliente [get]
func GetAllClients(c echo.Context) error {
	logger.Log.Info("[ClientController][GetAllClients] Intentando obtener clientes")
	clients, err := services.GetAllClients()
	if err != nil {
		logger.Log.Error("[ClientController][GetAllClients] Error al obtener clientes: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ClientController][GetAllClients] Clientes obtenidos: %d", len(clients))
	return helpers.RespondSuccess(c, "Clientes obtenidos", clients)
}

// @Summary Obtener cliente por ID
// @Description Devuelve los datos de un cliente específico.
// @Tags Clientes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Success 200 {object} dtos.Response{data=dtos.GetClientDto} "Cliente encontrado"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido"
// @Failure 404 {object} dtos.ErrorResponse "Cliente no encontrado"
// @Router /cliente/{id} [get]
func GetClientByID(c echo.Context) error {
	id := c.Param("id")
	logger.Log.Infof("[ClientController][GetClientByID] Intentando obtener cliente con ID: %s", id)
	ClientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ClientController][GetClientByID] Error al obtener cliente: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	client, err := services.GetClientByID(uint(ClientID))
	if err != nil {
		logger.Log.Error("[ClientController][GetClientByID] Error al obtener cliente: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ClientController][GetClientByID] Cliente obtenido con éxito: ID %d", ClientID)
	return helpers.RespondSuccess(c, "Cliente obtenido", client)
}

// @Summary Actualizar cliente
// @Description Actualiza los datos de un cliente específico.
// @Tags Clientes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Param request body dtos.ClientDTO true "Datos actualizados del cliente"
// @Success 200 {object} dtos.Response{data=nil} "Cliente actualizado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Datos o ID inválidos"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /cliente/{id} [put]
func UpdateClient(c echo.Context) error {
	id := c.Param("id")
	logger.Log.Infof("[ClientController][UpdateClient] Intentando actualizar cliente con ID: %s", id)
	ClientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ClientController][UpdateClient] Error al actualizar cliente: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	var clientDTO dtos.ClientDTO
	if err := c.Bind(&clientDTO); err != nil {
		logger.Log.Warn("[ClientController][UpdateClient] Error al actualizar cliente: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err = services.UpdateClient(uint(ClientID), clientDTO)
	if err != nil {
		logger.Log.Error("[ClientController][UpdateClient] Error al actualizar cliente: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ClientController][UpdateClient] Cliente actualizado con éxito: ID %d", ClientID)
	return helpers.RespondSuccess(c, "Cliente actualizado exitosamente", nil)
}

// @Summary Eliminar cliente
// @Description Elimina un cliente específico.
// @Tags Clientes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Success 200 {object} dtos.Response{data=nil} "Cliente eliminado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /cliente/{id} [delete]
func DeleteClient(c echo.Context) error {
	id := c.Param("id")
	logger.Log.Infof("[ClientController][DeleteClient] Intentando eliminar cliente con ID: %s", id)
	ClientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ClientController][DeleteClient] Error al eliminar cliente: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}
	err = services.DeleteClient(uint(ClientID))
	if err != nil {
		logger.Log.Error("[ClientController][DeleteClient] Error al eliminar cliente: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ClientController][DeleteClient] Cliente eliminado con éxito: ID %d", ClientID)
	return helpers.RespondSuccess(c, "Cliente eliminado exitosamente", nil)
}
