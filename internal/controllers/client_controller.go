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
	return helpers.RespondSuccess(c, "Cliente creado", nil)
}

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
	return helpers.RespondSuccess(c, "Cliente actualizado", nil)
}

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
	return helpers.RespondSuccess(c, "Cliente eliminado", nil)
}
