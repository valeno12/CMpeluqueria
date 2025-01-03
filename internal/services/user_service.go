package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"
	"peluqueria/middlewares"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(loginDto dtos.LoginDto) (dtos.LoginAnswerDto, error) {
	logger.Log.Infof("Intentando login para usuario: %s", loginDto.Username)

	// Validar datos
	if loginDto.Username == "" || loginDto.Password == "" {
		logger.Log.Warn("Login fallido: username o password faltante")
		return dtos.LoginAnswerDto{}, errors.New("username y password son obligatorios")
	}

	// Buscar usuario
	var user models.User
	if err := database.DB.Where("username = ?", loginDto.Username).First(&user).Error; err != nil {
		logger.Log.Warnf("Usuario no encontrado: %s", loginDto.Username)
		return dtos.LoginAnswerDto{}, errors.New("usuario no encontrado")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password)); err != nil {
		logger.Log.Warnf("Contraseña incorrecta para usuario: %s", loginDto.Username)
		return dtos.LoginAnswerDto{}, errors.New("contraseña incorrecta")
	}

	// Generar token
	jwtkey, jwterr := middlewares.GenerateToken(user)
	if jwterr != nil {
		logger.Log.Error("Error al generar el token para usuario: ", loginDto.Username, " - ", jwterr)
		return dtos.LoginAnswerDto{}, errors.New("error al generar el token")
	}

	return dtos.LoginAnswerDto{
		Username: user.Username,
		Token:    jwtkey,
	}, nil
}

func CreateUser(userDto dtos.UserDto) error {
	logger.Log.Infof("Intentando crear usuario: %s", userDto.Username)

	// Validar datos
	if userDto.Username == "" || userDto.Password == "" {
		logger.Log.Warn("Datos faltantes para crear usuario")
		return errors.New("username y password son obligatorios")
	}

	// Verificar si el usuario ya existe
	var existingUser models.User
	if err := database.DB.Where("username = ?", userDto.Username).First(&existingUser).Error; err == nil {
		logger.Log.Warnf("Usuario ya existente: %s", userDto.Username)
		return errors.New("el nombre de usuario ya está en uso")
	}

	// Encriptar contraseña
	userDto.Password = HashPassword(userDto.Password)

	// Crear el usuario
	user := models.User{
		Username: userDto.Username,
		Password: userDto.Password,
		RoleID:   userDto.RoleID,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		logger.Log.Error("Error al crear usuario: ", err)
		return err
	}

	logger.Log.Infof("Usuario creado con éxito: %s", userDto.Username)
	return nil
}

func GetAllUsers() ([]models.User, error) {
	logger.Log.Info("Obteniendo lista de todos los usuarios")

	var users []models.User
	result := database.DB.Preload("Role").Find(&users)
	if result.Error != nil {
		logger.Log.Error("Error al obtener usuarios: ", result.Error)
		return nil, result.Error
	}

	logger.Log.Infof("Usuarios obtenidos: %d", len(users))
	return users, nil
}

func GetUserByID(id uint) (models.User, error) {
	logger.Log.Infof("Buscando usuario con ID: %d", id)

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Log.Warnf("Usuario no encontrado: ID %d", id)
			return models.User{}, errors.New("usuario no encontrado")
		}
		logger.Log.Error("Error al buscar usuario: ", err)
		return models.User{}, err
	}

	logger.Log.Infof("Usuario encontrado: ID %d", id)
	return user, nil
}

func UpdateUser(userDto dtos.UserDto, id uint) error {
	logger.Log.Infof("Actualizando usuario con ID: %d", id)

	// Validar ID
	if id == 0 {
		logger.Log.Warn("ID del usuario faltante en actualización")
		return errors.New("el ID del usuario es obligatorio")
	}

	// Buscar usuario existente
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Log.Warnf("Usuario no encontrado para actualizar: ID %d", id)
			return errors.New("el usuario no existe")
		}
		logger.Log.Error("Error al buscar usuario para actualizar: ", err)
		return err
	}

	// Actualizar los campos necesarios
	if userDto.Username != "" {
		user.Username = userDto.Username
	}
	if userDto.Password != "" {
		user.Password = HashPassword(userDto.Password)
	}
	if userDto.RoleID != 0 {
		user.RoleID = userDto.RoleID
	}

	if err := database.DB.Save(&user).Error; err != nil {
		logger.Log.Error("Error al actualizar usuario: ", err)
		return errors.New("error al actualizar el usuario")
	}

	logger.Log.Infof("Usuario actualizado con éxito: ID %d", id)
	return nil
}

func DeleteUser(id uint) error {
	logger.Log.Infof("Eliminando usuario con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("ID del usuario faltante en eliminación")
		return errors.New("el ID del usuario es obligatorio")
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		logger.Log.Error("Error al eliminar usuario: ", err)
		return errors.New("error al eliminar usuario")
	}

	logger.Log.Infof("Usuario eliminado con éxito: ID %d", id)
	return nil
}

func HashPassword(pass string) string {
	logger.Log.Info("Encriptando contraseña")
	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), costo)
	if err != nil {
		logger.Log.Error("Error al encriptar contraseña: ", err)
	}
	return string(bytes)
}
