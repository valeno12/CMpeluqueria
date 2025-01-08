package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"
	"time"
)

func GetStockMovements(stockType string, month string) ([]dtos.StockMovementDto, error) {
	logger.Log.Info("[StockService][GetStockMovements] Obteniendo movimientos de stock")

	var movements []models.StockMovement
	query := database.DB.Unscoped().Preload("Product")

	// Filtrar por tipo de movimiento
	if stockType == "entry" {
		logger.Log.Info("[StockService][GetStockMovements] Filtrando por movimientos de entrada")
		query = query.Where("quantity > 0")
	} else if stockType == "exit" {
		logger.Log.Info("[StockService][GetStockMovements] Filtrando por movimientos de salida")
		query = query.Where("quantity < 0")
	}

	// Filtrar por mes
	if month != "" {
		logger.Log.Infof("[StockService][GetStockMovements] Aplicando filtro por mes: %s", month)
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

	// Mapear resultados a DTO
	var movementDtos []dtos.StockMovementDto
	for _, movement := range movements {
		movementDtos = append(movementDtos, dtos.StockMovementDto{
			ID:             movement.ID,
			ProductID:      movement.ProductID,
			ProductName:    movement.Product.Name,
			ProductBrand:   movement.Product.Brand,
			ProductUnit:    movement.ProductUnit,
			Quantity:       movement.Quantity,
			PackageCount:   movement.PackageCount,
			UnitPerPackage: movement.UnitPerPackage,
			UnityPrice:     movement.UnityPrice,
			Reason:         movement.Reason,
			CreatedAt:      movement.CreatedAt.Format("02/01/2006 15:04"),
		})
	}

	logger.Log.Infof("[StockService][GetStockMovements] Movimientos obtenidos: %d", len(movementDtos))
	return movementDtos, nil
}

func GetStockMovementsByProduct(productID uint, stockType string, month string) ([]dtos.StockMovementDto, error) {
	logger.Log.Infof("[StockService][GetStockMovementsByProduct] Obteniendo movimientos de stock para producto ID: %d", productID)

	var movements []models.StockMovement
	query := database.DB.Unscoped().Where("product_id = ?", productID).Preload("Product")

	// Filtrar por tipo de movimiento
	if stockType == "entry" {
		logger.Log.Info("[StockService][GetStockMovementsByProduct] Filtrando por movimientos de entrada")
		query = query.Where("quantity > 0")
	} else if stockType == "exit" {
		logger.Log.Info("[StockService][GetStockMovementsByProduct] Filtrando por movimientos de salida")
		query = query.Where("quantity < 0")
	}

	// Filtrar por mes
	if month != "" {
		logger.Log.Infof("[StockService][GetStockMovementsByProduct] Aplicando filtro por mes: %s", month)
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

	// Mapear resultados a DTO
	var movementDtos []dtos.StockMovementDto
	for _, movement := range movements {
		movementDtos = append(movementDtos, dtos.StockMovementDto{
			ID:             movement.ID,
			ProductID:      movement.ProductID,
			ProductName:    movement.Product.Name,
			ProductBrand:   movement.Product.Brand,
			ProductUnit:    movement.ProductUnit,
			Quantity:       movement.Quantity,
			PackageCount:   movement.PackageCount,
			UnitPerPackage: movement.UnitPerPackage,
			UnityPrice:     movement.UnityPrice,
			Reason:         movement.Reason,
			CreatedAt:      movement.CreatedAt.Format("02/01/2006 15:04"),
		})
	}

	logger.Log.Infof("[StockService][GetStockMovementsByProduct] Movimientos obtenidos: %d", len(movementDtos))
	return movementDtos, nil
}

func parseMonthFilter(month string) (time.Time, time.Time, error) {
	logger.Log.Infof("[StockService][parseMonthFilter] Parseando filtro de mes: %s", month)

	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		logger.Log.Warn("[StockService][parseMonthFilter] Error al parsear filtro de mes: ", err)
		return time.Time{}, time.Time{}, errors.New("formato de mes inválido. Use 'YYYY-MM'")
	}

	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	logger.Log.Infof("[StockService][parseMonthFilter] Filtro parseado correctamente - Inicio: %s, Fin: %s", startDate, endDate)
	return startDate, endDate, nil
}
