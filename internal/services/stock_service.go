package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/models"
	"peluqueria/logger"
	"time"
)

func GetStockMovements(stockType string, month string) ([]models.StockMovement, error) {
	logger.Log.Info("[StockService][GetStockMovements] Obteniendo movimientos de stock")

	var movements []models.StockMovement
	query := database.DB.Preload("Product")

	// Filtrar por tipo de movimiento
	if stockType == "entry" {
		query = query.Where("quantity > 0")
	} else if stockType == "exit" {
		query = query.Where("quantity < 0")
	}

	// Filtrar por mes
	if month != "" {
		startDate, endDate, err := parseMonthFilter(month)
		if err != nil {
			logger.Log.Warn("[StockService][GetStockMovements] Filtro de mes inválido: ", err)
			return nil, err
		}
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	// Ejecutar consulta
	if err := query.Find(&movements).Error; err != nil {
		logger.Log.Error("[StockService][GetStockMovements] Error al obtener movimientos: ", err)
		return nil, errors.New("error al obtener movimientos de stock")
	}

	return movements, nil
}

func GetStockMovementsByProduct(productID uint, stockType string, month string) ([]models.StockMovement, error) {
	logger.Log.Infof("[StockService][GetStockMovementsByProduct] Obteniendo movimientos de stock para producto ID: %d", productID)

	var movements []models.StockMovement
	query := database.DB.Where("product_id = ?", productID).Preload("Product")

	// Filtrar por tipo de movimiento
	if stockType == "entry" {
		query = query.Where("quantity > 0")
	} else if stockType == "exit" {
		query = query.Where("quantity < 0")
	}

	// Filtrar por mes
	if month != "" {
		startDate, endDate, err := parseMonthFilter(month)
		if err != nil {
			logger.Log.Warn("[StockService][GetStockMovementsByProduct] Filtro de mes inválido: ", err)
			return nil, err
		}
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	// Ejecutar consulta
	if err := query.Find(&movements).Error; err != nil {
		logger.Log.Error("[StockService][GetStockMovementsByProduct] Error al obtener movimientos: ", err)
		return nil, errors.New("error al obtener movimientos de stock por producto")
	}

	return movements, nil
}

func parseMonthFilter(month string) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("formato de mes inválido. Use 'YYYY-MM'")
	}
	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return startDate, endDate, nil
}
