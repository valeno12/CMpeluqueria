package helpers

import (
	"errors"
	"net/http"
	"peluqueria/internal/dtos"
	"time"

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

func ParseCustomDate(dateString string) (time.Time, error) {
	// Formato esperado
	layout := "15:04 02/01/2006"

	// Cargar la zona horaria
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return time.Time{}, errors.New("no se pudo cargar la zona horaria")
	}

	// Parsear la fecha con la zona horaria especificada
	parsedTime, err := time.ParseInLocation(layout, dateString, loc)
	if err != nil {
		return time.Time{}, errors.New("formato de fecha inválido, debe ser HH:MM DD/MM/YYYY")
	}

	return parsedTime, nil
}
