package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/logger"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware valida el token JWT y añade información al contexto.
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener el token desde el header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.Warn("Token de autorización no encontrado")
			return respondError(c, http.StatusUnauthorized, "Token de autorización no encontrado")
		}

		// Verificar formato Bearer
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Log.Warn("Formato de token inválido: ", authHeader)
			return respondError(c, http.StatusUnauthorized, "Formato de token inválido")
		}

		tokenString := tokenParts[1]

		// Cargar la clave secreta
		secret := []byte(os.Getenv("SECRET_JWT"))
		if len(secret) == 0 {
			logger.Log.Error("Clave secreta para JWT no encontrada")
			return respondError(c, http.StatusInternalServerError, "Error interno del servidor")
		}

		// Verificar y decodificar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar que el algoritmo usado es HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Log.Warn("Método de firma inválido")
				return nil, fmt.Errorf("método de firma inválido")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			logger.Log.Warn("Token inválido o expirado: ", err)
			return respondError(c, http.StatusUnauthorized, "Token inválido")
		}

		// Extraer claims del token
		claims := token.Claims.(jwt.MapClaims)

		// Validar y convertir claims
		userID, ok := claims["sub"].(float64)
		if !ok {
			logger.Log.Warn("Token inválido: 'sub' no encontrado o formato incorrecto")
			return respondError(c, http.StatusUnauthorized, "Token inválido: sub no encontrado o formato incorrecto")
		}

		roleID, ok := claims["role"].(float64)
		if !ok {
			logger.Log.Warn("Token inválido: 'role' no encontrado o formato incorrecto")
			return respondError(c, http.StatusUnauthorized, "Token inválido: role no encontrado o formato incorrecto")
		}

		logger.Log.Infof("Token válido para UserID: %d, RoleID: %d", uint(userID), uint(roleID))

		// Añadir información al contexto
		c.Set("user_id", uint(userID))
		c.Set("role_id", uint(roleID))

		// Continuar con el siguiente handler
		return next(c)
	}
}

// PermissionMiddleware verifica si el rol del usuario tiene el permiso requerido.
func PermissionMiddleware(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleID := c.Get("role_id").(uint)

			// Verificar si el rol tiene el permiso
			allowed, err := hasPermission(roleID, permission)
			if err != nil {
				logger.Log.Error("Error al verificar permisos para RoleID: ", roleID, ", Permission: ", permission, ", Error: ", err)
				return respondError(c, http.StatusInternalServerError, "Error al verificar permisos")
			}

			if !allowed {
				logger.Log.Warn("Permiso denegado para RoleID: ", roleID, ", Permission: ", permission)
				return respondError(c, http.StatusForbidden, "No tienes permiso para realizar esta acción")
			}

			logger.Log.Infof("Permiso concedido para RoleID: %d, Permission: %s", roleID, permission)
			return next(c)
		}
	}
}

// hasPermission verifica si un rol tiene un permiso específico.
func hasPermission(roleID uint, permission string) (bool, error) {
	var count int64
	err := database.DB.Table("role_permissions").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, permission).
		Count(&count).Error

	if err != nil {
		logger.Log.Error("Error al consultar permisos para RoleID: ", roleID, ", Permission: ", permission, ", Error: ", err)
		return false, err
	}

	return count > 0, nil
}

func RouteLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Loguear detalles de la solicitud
		logger.Log.Infof("Ruta accedida: método=%s, ruta=%s", c.Request().Method, c.Request().URL.Path)

		// Continuar con el siguiente middleware o handler
		return next(c)
	}
}

// Genera una respuesta de error
func respondError(c echo.Context, status int, message string) error {
	logger.Log.Warnf("RespondError - Status: %d, Message: %s", status, message)
	return c.JSON(status, dtos.Response{
		Status:  "error",
		Message: message,
	})
}
