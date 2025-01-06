package controllers

import (
	"net/http"
	"peluqueria/internal/helpers"
	"peluqueria/internal/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetStockMovements(c echo.Context) error {
	stockType := c.QueryParam("type")
	month := c.QueryParam("month")

	movements, err := services.GetStockMovements(stockType, month)
	if err != nil {
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.RespondSuccess(c, "Movimientos de stock obtenidos", movements)
}

func GetStockMovementsByProduct(c echo.Context) error {
	// Obtener ID del producto
	productIDParam := c.Param("id")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	// Obtener filtros
	stockType := c.QueryParam("type")
	month := c.QueryParam("month")

	movements, err := services.GetStockMovementsByProduct(uint(productID), stockType, month)
	if err != nil {
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.RespondSuccess(c, "Movimientos de stock obtenidos por producto", movements)
}
