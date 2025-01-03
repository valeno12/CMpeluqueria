package middlewares

import (
	"os"
	"peluqueria/internal/models"
	"peluqueria/logger"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user models.User) (string, error) {
	// Obtener la clave secreta
	key := []byte(os.Getenv("SECRET_JWT"))
	if len(key) == 0 {
		logger.Log.Error("Clave secreta para JWT no encontrada en las variables de entorno")
		return "", jwt.ErrSignatureInvalid
	}

	// Log de generación de token
	logger.Log.Infof("Generando token para el usuario ID: %d, RoleID: %d", user.ID, user.RoleID)

	// Crear el token con claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.RoleID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// Firmar el token
	tokenString, err := token.SignedString(key)
	if err != nil {
		logger.Log.Error("Error al firmar el token JWT: ", err)
		return "", err
	}

	// Log de éxito
	logger.Log.Infof("Token generado exitosamente para el usuario ID: %d", user.ID)
	return tokenString, nil
}
