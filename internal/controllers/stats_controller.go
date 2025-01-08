package controllers

import (
	"net/http"
	"peluqueria/internal/helpers"
	"peluqueria/internal/services"

	"github.com/labstack/echo/v4"
)

func GetMonthlyStatistics(c echo.Context) error {
	month := c.QueryParam("month")
	if month == "" {
		return helpers.RespondError(c, http.StatusBadRequest, "El parámetro 'month' es obligatorio")
	}

	statistics, err := services.GetMonthlyStatistics(month)
	if err != nil {
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.RespondSuccess(c, "Estadísticas mensuales obtenidas", statistics)
}
