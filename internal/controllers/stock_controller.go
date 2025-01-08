package controllers

import (
	"net/http"
	"peluqueria/internal/helpers"
	"peluqueria/internal/services"
	"peluqueria/logger"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetStockMovements obtiene todos los movimientos de stock.
// @Summary Obtener movimientos de stock
// @Description Devuelve una lista de todos los movimientos de stock, filtrados opcionalmente por tipo y mes.
// @Tags Movimientos de Stock
// @Accept json
// @Produce json
// @Param type query string false "Filtrar por tipo de movimiento ('entry' para entradas, 'exit' para salidas)" example:"entry"
// @Param month query string false "Filtrar por mes (formato YYYY-MM)" example:"2025-01"
// @Success 200 {object} dtos.Response{data=[]dtos.StockMovementDto} "Movimientos de stock obtenidos"
// @Failure 500 {object} dtos.Response{data=nil} "Error al obtener movimientos de stock"
// @Router /stock-movements [get]
// @Security BearerAuth
func GetStockMovements(c echo.Context) error {
	logger.Log.Info("[StockController][GetStockMovements] Obteniendo movimientos de stock")

	stockType := c.QueryParam("type")
	month := c.QueryParam("month")

	logger.Log.Infof("[StockController][GetStockMovements] Parámetros recibidos - Type: %s, Month: %s", stockType, month)

	movements, err := services.GetStockMovements(stockType, month)
	if err != nil {
		logger.Log.Error("[StockController][GetStockMovements] Error al obtener movimientos de stock: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "Error al obtener movimientos de stock")
	}

	logger.Log.Infof("[StockController][GetStockMovements] Movimientos obtenidos: %d", len(movements))
	return helpers.RespondSuccess(c, "Movimientos de stock obtenidos", movements)
}

// GetStockMovementsByProduct obtiene los movimientos de stock para un producto específico.
// @Summary Obtener movimientos de stock por producto
// @Description Devuelve los movimientos de stock asociados a un producto específico, filtrados opcionalmente por tipo y mes.
// @Tags Movimientos de Stock
// @Accept json
// @Produce json
// @Success 200 {object} dtos.Response{data=[]dtos.StockMovementDto} "Movimientos de stock obtenidos por producto"
// @Failure 400 {object} dtos.Response{data=nil} "ID del producto inválido"
// @Failure 500 {object} dtos.Response{data=nil} "Error al obtener movimientos de stock por producto"
// @Router /stock-movements/product/{id} [get]
// @Security BearerAuth
func GetStockMovementsByProduct(c echo.Context) error {
	logger.Log.Info("[StockController][GetStockMovementsByProduct] Obteniendo movimientos de stock por producto")

	productIDParam := c.Param("id")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		logger.Log.Warn("[StockController][GetStockMovementsByProduct] ID del producto inválido: ", productIDParam)
		return helpers.RespondError(c, http.StatusBadRequest, "ID del producto inválido")
	}

	stockType := c.QueryParam("type")
	month := c.QueryParam("month")

	logger.Log.Infof("[StockController][GetStockMovementsByProduct] Parámetros recibidos - ID: %d, Type: %s, Month: %s", productID, stockType, month)

	movements, err := services.GetStockMovementsByProduct(uint(productID), stockType, month)
	if err != nil {
		logger.Log.Error("[StockController][GetStockMovementsByProduct] Error al obtener movimientos de stock por producto: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, "Error al obtener movimientos de stock por producto")
	}

	logger.Log.Infof("[StockController][GetStockMovementsByProduct] Movimientos obtenidos: %d", len(movements))
	return helpers.RespondSuccess(c, "Movimientos de stock obtenidos por producto", movements)
}
