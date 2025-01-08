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

// @Summary Iniciar sesión
// @Description Permite a un usuario autenticarse en el sistema.
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body dtos.LoginDto true "Datos de inicio de sesión"
// @Success 200 {object} dtos.Response{data=dtos.LoginAnswerDto} "Token de acceso"
// @Failure 400 {object} dtos.ErrorResponse "Datos inválidos. Ejemplo: El formato de los datos es incorrecto"
// @Failure 401 {object} dtos.ErrorResponse "Credenciales inválidas. Ejemplo: Usuario o contraseña incorrectos"
// @Router /login [post]
func Login(c echo.Context) error {
	logger.Log.Info("[UserController][Login] Intentando iniciar sesión")
	var login dtos.LoginDto
	if err := c.Bind(&login); err != nil {
		logger.Log.Warn("[UserController][Login] Login fallido: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "El formato de los datos es incorrecto")
	}

	token, err := services.Login(login)
	if err != nil {
		logger.Log.Warn("[LoginController][Login] Login fallido: ", err)
		return helpers.RespondError(c, http.StatusUnauthorized, "Usuario o contraseña incorrectos")
	}

	logger.Log.Infof("[LoginController][Login] Login exitoso para usuario: %s", login.Username)
	return helpers.RespondSuccess(c, "Login exitoso", token)
}

// @Summary Crear usuario
// @Description Crea un nuevo usuario en el sistema.
// @Tags Usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body dtos.UserDto true "Datos del usuario"
// @Success 200 {object} dtos.Response{data=string} "Usuario creado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Datos inválidos. Ejemplo: El nombre de usuario ya existe"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor. Ejemplo: No se pudo conectar a la base de datos"
// @Router /usuarios [post]
func CreateUser(c echo.Context) error {
	logger.Log.Info("[UserController][CreateUser] Intentando crear usuario")
	var user dtos.UserDto
	if err := c.Bind(&user); err != nil {
		logger.Log.Warn("[UserController][CreateUser] Error al crear usuario: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "El nombre de usuario ya existe")
	}

	err := services.CreateUser(user)
	if err != nil {
		logger.Log.Error("[UserController][CreateUser] Error al crear usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo conectar a la base de datos")
	}

	logger.Log.Infof("[UserController][CreateUser] Usuario creado: %s", user.Username)
	return helpers.RespondSuccess(c, "Usuario creado exitosamente", nil)
}

// @Summary Obtener todos los usuarios
// @Description Devuelve una lista de todos los usuarios registrados.
// @Tags Usuarios
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.Response{data=[]dtos.GetUserDto} "Usuarios encontrados"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor. Ejemplo: No se pudo obtener la lista de usuarios"
// @Router /usuarios [get]
func GetAllUsers(c echo.Context) error {
	logger.Log.Info("[UserController][GetAllUsers] Obteniendo todos los usuarios")
	users, err := services.GetAllUsers()
	if err != nil {
		logger.Log.Error("[UserController][GetAllUsers] Error al obtener usuarios: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo obtener la lista de usuarios")
	}

	logger.Log.Infof("[UserController][GetAllUsers] Usuarios obtenidos: %d usuarios encontrados", len(users))
	return helpers.RespondSuccess(c, "Usuarios encontrados", users)
}

// @Summary Obtener usuario por ID
// @Description Devuelve los datos de un usuario específico.
// @Tags Usuarios
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Success 200 {object} dtos.Response{data=dtos.GetUserDto} "Usuario encontrado"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido. Ejemplo: El ID proporcionado no es válido"
// @Failure 404 {object} dtos.ErrorResponse "Usuario no encontrado. Ejemplo: El usuario solicitado no existe"
// @Router /usuarios/{id} [get]
func GetUserByID(c echo.Context) error {
	userIDParam := c.Param("id")
	logger.Log.Infof("[UserController][GetUserByID] Obteniendo usuario con ID: %s", userIDParam)

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[UserController][GetUserByID] ID de usuario inválido: ", userIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "El ID proporcionado no es válido")
	}

	user, err := services.GetUserByID(uint(userID))
	if err != nil {
		logger.Log.Warnf("[UserController][GetUserByID] Usuario con ID %d no encontrado", userID)
		return helpers.RespondError(c, http.StatusNotFound, "El usuario solicitado no existe")
	}

	logger.Log.Infof("[UserController][GetUserByID] Usuario encontrado: ID %d", userID)
	return helpers.RespondSuccess(c, "Usuario encontrado", user)
}

// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario específico.
// @Tags Usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Param user body dtos.UserDto true "Datos del usuario"
// @Success 200 {object} dtos.Response{data=string} "Usuario actualizado correctamente"
// @Failure 400 {object} dtos.ErrorResponse "Datos o ID inválidos. Ejemplo: Los datos enviados son inválidos"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor. Ejemplo: No se pudo actualizar el usuario"
// @Router /usuarios/{id} [put]
func UpdateUser(c echo.Context) error {
	userIDParam := c.Param("id")
	logger.Log.Infof("[UserController][UpdateUser] Intentando actualizar usuario con ID: %s", userIDParam)

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[UserController][UpdateUser] ID de usuario inválido: ", userIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "El ID proporcionado no es válido")
	}

	var user dtos.UserDto
	if err := c.Bind(&user); err != nil {
		logger.Log.Warn("[UserController][UpdateUser] Error al actualizar usuario: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Los datos enviados son inválidos")
	}

	err = services.UpdateUser(user, uint(userID))
	if err != nil {
		logger.Log.Error("[UserController][UpdateUser] Error al actualizar usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo actualizar el usuario")
	}

	logger.Log.Infof("[UserController][UpdateUser] Usuario actualizado: ID %d", userID)
	return helpers.RespondSuccess(c, "Usuario actualizado correctamente", nil)
}

// @Summary Eliminar usuario
// @Description Elimina un usuario específico.
// @Tags Usuarios
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Success 200 {object} dtos.Response{data=string} "Usuario eliminado correctamente"
// @Failure 400 {object} dtos.ErrorResponse "ID inválido. Ejemplo: El ID proporcionado no es válido"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor. Ejemplo: No se pudo eliminar el usuario"
// @Router /usuarios/{id} [delete]
func DeleteUser(c echo.Context) error {
	userIDParam := c.Param("id")
	logger.Log.Infof("[UserController][DeleteUser] Intentando eliminar usuario con ID: %s", userIDParam)

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[UserController][DeleteUser] ID de usuario inválido: ", userIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "El ID proporcionado no es válido")
	}

	err = services.DeleteUser(uint(userID))
	if err != nil {
		logger.Log.Error("[UserController][DeleteUser] Error al eliminar usuario: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "No se pudo eliminar el usuario")
	}

	logger.Log.Infof("[UserController][DeleteUser] Usuario eliminado: ID %d", userID)
	return helpers.RespondSuccess(c, "Usuario eliminado correctamente", nil)
}
