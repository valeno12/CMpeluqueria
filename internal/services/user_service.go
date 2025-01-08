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
	logger.Log.Infof("[UserService][Login] Intentando login para usuario: %s", loginDto.Username)

	if loginDto.Username == "" || loginDto.Password == "" {
		logger.Log.Warn("[UserService][Login] Username o password faltante")
		return dtos.LoginAnswerDto{}, errors.New("username y password son obligatorios")
	}

	var user models.User
	if err := database.DB.Where("username = ?", loginDto.Username).First(&user).Error; err != nil {
		logger.Log.Warnf("[UserService][Login] Usuario no encontrado: %s", loginDto.Username)
		return dtos.LoginAnswerDto{}, errors.New("usuario no encontrado")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password)); err != nil {
		logger.Log.Warnf("[UserService][Login] Contraseña incorrecta para usuario: %s", loginDto.Username)
		return dtos.LoginAnswerDto{}, errors.New("contraseña incorrecta")
	}

	jwtKey, jwtErr := middlewares.GenerateToken(user)
	if jwtErr != nil {
		logger.Log.Error("[UserService][Login] Error al generar el token: ", jwtErr)
		return dtos.LoginAnswerDto{}, errors.New("error al generar el token")
	}

	logger.Log.Infof("[UserService][Login] Login exitoso para usuario: %s", loginDto.Username)
	return dtos.LoginAnswerDto{Username: user.Username, Token: jwtKey}, nil
}

func CreateUser(userDto dtos.UserDto) error {
	logger.Log.Infof("[UserService][CreateUser] Intentando crear usuario: %s", userDto.Username)

	if userDto.Username == "" || userDto.Password == "" {
		logger.Log.Warn("[UserService][CreateUser] Username o password faltante")
		return errors.New("username y password son obligatorios")
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", userDto.Username).First(&existingUser).Error; err == nil {
		logger.Log.Warnf("[UserService][CreateUser] Usuario ya existente: %s", userDto.Username)
		return errors.New("el nombre de usuario ya está en uso")
	}

	if userDto.RoleID != 0 {
		var role models.Role
		if err := database.DB.First(&role, userDto.RoleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Log.Warnf("[UserService][CreateUser] Rol no encontrado: ID %d", userDto.RoleID)
				return errors.New("el rol asignado no existe")
			}
			logger.Log.Error("[UserService][CreateUser] Error al buscar rol: ", err)
			return errors.New("error al validar el rol asignado")
		}
	}

	// Hashear la contraseña
	userDto.Password = HashPassword(userDto.Password)
	user := models.User{
		Username: userDto.Username,
		Password: userDto.Password,
		RoleID:   userDto.RoleID,
	}

	// Crear el usuario
	if err := database.DB.Create(&user).Error; err != nil {
		logger.Log.Error("[UserService][CreateUser] Error al crear usuario: ", err)
		return err
	}

	logger.Log.Infof("[UserService][CreateUser] Usuario creado con éxito: %s", userDto.Username)
	return nil
}

func GetAllUsers() ([]dtos.GetUserDto, error) {
	logger.Log.Info("[UserService][GetAllUsers] Obteniendo todos los usuarios")

	var users []models.User
	if err := database.DB.Preload("Role").Find(&users).Error; err != nil {
		logger.Log.Error("[UserService][GetAllUsers] Error al obtener usuarios: ", err)
		return nil, err
	}

	var userDtos []dtos.GetUserDto
	for _, user := range users {
		roleName := ""
		if user.Role.ID != 0 {
			roleName = user.Role.Name
		}
		userDtos = append(userDtos, dtos.GetUserDto{
			ID:       user.ID,
			Username: user.Username,
			RoleName: roleName,
		})
	}

	logger.Log.Infof("[UserService][GetAllUsers] Usuarios obtenidos: %d", len(users))
	return userDtos, nil
}

func GetUserByID(id uint) (dtos.GetUserDto, error) {
	logger.Log.Infof("[UserService][GetUserByID] Buscando usuario con ID: %d", id)

	var user models.User
	if err := database.DB.Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[UserService][GetUserByID] Usuario no encontrado: ID %d", id)
			return dtos.GetUserDto{}, errors.New("usuario no encontrado")
		}
		logger.Log.Error("[UserService][GetUserByID] Error al buscar usuario: ", err)
		return dtos.GetUserDto{}, err
	}

	roleName := ""
	if user.Role.ID != 0 {
		roleName = user.Role.Name
	}

	logger.Log.Infof("[UserService][GetUserByID] Usuario encontrado: ID %d", id)
	return dtos.GetUserDto{
		ID:       user.ID,
		Username: user.Username,
		RoleName: roleName,
	}, nil
}

func UpdateUser(userDto dtos.UserDto, id uint) error {
	logger.Log.Infof("[UserService][UpdateUser] Actualizando usuario con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("[UserService][UpdateUser] ID del usuario faltante")
		return errors.New("el ID del usuario es obligatorio")
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[UserService][UpdateUser] Usuario no encontrado: ID %d", id)
			return errors.New("el usuario no existe")
		}
		logger.Log.Error("[UserService][UpdateUser] Error al buscar usuario: ", err)
		return err
	}

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
		logger.Log.Error("[UserService][UpdateUser] Error al actualizar usuario: ", err)
		return errors.New("error al actualizar usuario")
	}

	logger.Log.Infof("[UserService][UpdateUser] Usuario actualizado con éxito: ID %d", id)
	return nil
}

func DeleteUser(id uint) error {
	logger.Log.Infof("[UserService][DeleteUser] Eliminando usuario con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("[UserService][DeleteUser] ID del usuario faltante")
		return errors.New("el ID del usuario es obligatorio")
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		logger.Log.Error("[UserService][DeleteUser] Error al eliminar usuario: ", err)
		return errors.New("error al eliminar usuario")
	}

	logger.Log.Infof("[UserService][DeleteUser] Usuario eliminado con éxito: ID %d", id)
	return nil
}

func HashPassword(pass string) string {
	logger.Log.Info("[UserService][HashPassword] Encriptando contraseña")
	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), costo)
	if err != nil {
		logger.Log.Error("[UserService][HashPassword] Error al encriptar contraseña: ", err)
	}
	return string(bytes)
}
