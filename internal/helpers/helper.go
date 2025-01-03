package helpers

import (
	"net/http"
	"peluqueria/internal/dtos"

	"github.com/labstack/echo/v4"
)

// Genera una respuesta exitosa
func RespondSuccess(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, dtos.Response{
		Status:  "ok",
		Message: message,
		Data:    data,
	})
}

// Genera una respuesta de error
func RespondError(c echo.Context, status int, message string) error {
	return c.JSON(status, dtos.Response{
		Status:  "error",
		Message: message,
	})
}
