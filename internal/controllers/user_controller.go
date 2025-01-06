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

func Login(c echo.Context) error {
	logger.Log.Info("[UserController][Login] Intentando iniciar sesión")
	var login dtos.LoginDto
	if err := c.Bind(&login); err != nil {
		logger.Log.Warn("[UserController][Login] Login fallido: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	token, err := services.Login(login)
	if err != nil {
		logger.Log.Warn("[LoginController][Login] Login fallido: ", err)
		return helpers.RespondError(c, http.StatusUnauthorized, err.Error())
	}

	logger.Log.Infof("[LoginController][Login] Login exitoso para usuario: %s", login.Username)
	return helpers.RespondSuccess(c, "Login exitoso", token)
}

func CreateUser(c echo.Context) error {
	logger.Log.Info("[UserController][CreateUser] Intentando crear usuario")
	var user dtos.UserDto
	if err := c.Bind(&user); err != nil {
		logger.Log.Warn("[UserController][CreateUser] Error al crear usuario: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err := services.CreateUser(user)
	if err != nil {
		logger.Log.Error("[UserController][CreateUser] Error al crear usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[UserController][CreateUser] Usuario creado: %s", user.Username)
	return helpers.RespondSuccess(c, "Usuario creado", nil)
}

func GetAllUsers(c echo.Context) error {
	logger.Log.Info("[UserController][GetAllUsers] Obteniendo todos los usuarios")
	users, err := services.GetAllUsers()
	if err != nil {
		logger.Log.Error("[UserController][GetAllUsers] Error al obtener usuarios: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[UserController][GetAllUsers] Usuarios obtenidos: %d usuarios encontrados", len(users))
	return helpers.RespondSuccess(c, "Usuarios encontrados", users)
}

func GetUserByID(c echo.Context) error {
	userIDParam := c.Param("id")
	logger.Log.Infof("[UserController][GetUserByID] Obteniendo usuario con ID: %s", userIDParam)

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[UserController][GetUserByID] ID de usuario inválido: ", userIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	user, err := services.GetUserByID(uint(userID))
	if err != nil {
		logger.Log.Warnf("[UserController][GetUserByID] Usuario con ID %d no encontrado", userID)
		return helpers.RespondError(c, http.StatusNotFound, err.Error())
	}

	logger.Log.Infof("[UserController][GetUserByID] Usuario encontrado: ID %d", userID)
	return helpers.RespondSuccess(c, "Usuario encontrado", user)
}

func UpdateUser(c echo.Context) error {
	userIDParam := c.Param("id")
	logger.Log.Infof("[UserController][UptadeUser] Intentando actualizar usuario con ID: %s", userIDParam)

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[UserController][UptadeUser] ID de usuario inválido: ", userIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var user dtos.UserDto
	if err := c.Bind(&user); err != nil {
		logger.Log.Warn("[UserController][UptadeUser] Error al actualizar usuario: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.UpdateUser(user, uint(userID))
	if err != nil {
		logger.Log.Error("[UserController][UptadeUser] Error al actualizar usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[UserController][UptadeUser] Usuario actualizado: ID %d", userID)
	return helpers.RespondSuccess(c, "Usuario actualizado correctamente", nil)
}

func DeleteUser(c echo.Context) error {
	userIDParam := c.Param("id")
	logger.Log.Infof("[UserController][DeleteUser] Intentando eliminar usuario con ID: %s", userIDParam)

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[UserController][DeleteUser] ID de usuario inválido: ", userIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	err = services.DeleteUser(uint(userID))
	if err != nil {
		logger.Log.Error("[UserController][DeleteUser] Error al eliminar usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Infof("[UserController][DeleteUser] Usuario eliminado: ID %d", userID)
	return helpers.RespondSuccess(c, "Usuario eliminado correctamente", nil)
}
