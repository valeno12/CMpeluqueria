package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"

	"gorm.io/gorm"
)

func CreateProduct(ProductDto dtos.CreateProductDto) error {
	logger.Log.Infof("[ProductService][CreateProduct] Intentando crear producto: %s", ProductDto.Name)
	if ProductDto.Name == "" || ProductDto.Unit == "" || ProductDto.Brand == "" {
		logger.Log.Warn("[ProductService][CreateProduct] Datos faltantes para crear producto")
		return errors.New("nombre, unidad y marca son obligatorios")
	}
	if ProductDto.PackageCount <= 0 || ProductDto.UnitPerPackage <= 0 {
		logger.Log.Warn("[ProductService][CreateProduct] PackageCount y UnitPerPackage deben ser mayores a 0")
		return errors.New("cantidad de paquetes y unidades por paquete deben ser mayores a 0")
	}
	var existingProduct models.Product
	if err := database.DB.Where("name = ? AND brand = ?", ProductDto.Name, ProductDto.Brand).First(&existingProduct).Error; err == nil {
		logger.Log.Warnf("[ProductService][CreateProduct] Producto ya existe: %s %s", ProductDto.Name, ProductDto.Brand)
		return errors.New("producto ya existe")
	}

	quantity := ProductDto.PackageCount * ProductDto.UnitPerPackage
	product := models.Product{
		Name:          ProductDto.Name,
		Unit:          ProductDto.Unit,
		Brand:         ProductDto.Brand,
		Quantity:      quantity,
		LowStockAlert: ProductDto.LowStockAlert,
	}

	movement := models.StockMovement{
		ProductID:  product.ID,
		Quantity:   quantity,
		Reason:     "Inventario inicial",
		UnityPrice: ProductDto.UnityPrice,
		TotalPrice: ProductDto.UnityPrice * ProductDto.PackageCount,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}
		movement.ProductID = product.ID
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Log.Error("[ProductService][CreateProduct] Error en transacción: ", err)
		return errors.New("error al crear producto y registrar movimiento inicial" + err.Error())
	}

	return nil
}

func GetAllProducts() ([]models.Product, error) {
	logger.Log.Info("[ProductService][GetAllProducts] Intentando obtener productos")

	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		logger.Log.Error("[ProductService][GetAllProducts] Error al obtener productos: ", err)
		return nil, errors.New("error al obtener productos")
	}

	return products, nil
}

func GetProductByID(id uint) (models.Product, error) {
	logger.Log.Infof("[ProductService][GetProductByID] Intentando obtener producto con ID: %d", id)

	var product models.Product
	if err := database.DB.Where("id = ?", id).First(&product).Error; err != nil {
		logger.Log.Error("[ProductService][GetProductByID] Error al obtener producto: ", err)
		return models.Product{}, errors.New("error al obtener producto")
	}

	return product, nil
}

func UpdateProduct(id uint, ProductDto dtos.UpdateProductDto) error {
	logger.Log.Infof("[ProductService][UpdateProduct] Intentando actualizar producto con ID: %d", id)

	product, err := GetProductByID(id)
	if err != nil {
		return err
	}

	// Validar y asignar valores
	if ProductDto.Name != "" {
		product.Name = ProductDto.Name
	}
	if ProductDto.Unit != "" {
		product.Unit = ProductDto.Unit
	}
	if ProductDto.Brand != "" {
		product.Brand = ProductDto.Brand
	}
	if ProductDto.LowStockAlert > 0 {
		product.LowStockAlert = ProductDto.LowStockAlert
	}

	if err := database.DB.Save(&product).Error; err != nil {
		logger.Log.Error("[ProductService][UpdateProduct] Error al actualizar producto: ", err)
		return errors.New("error al actualizar producto")
	}

	logger.Log.Infof("[ProductService][UpdateProduct] Producto actualizado con éxito: ID %d", id)
	return nil
}

func DeleteProduct(id uint) error {
	logger.Log.Infof("[ProductService][DeleteProduct] Intentando eliminar producto con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("[ProductService][DeleteProduct] ID del producto faltante en eliminación")
		return errors.New("el ID del producto es obligatorio")
	}

	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		logger.Log.Error("[ProductService][DeleteProduct] Error al eliminar producto: ", err)
		return errors.New("error al eliminar producto")
	}

	logger.Log.Infof("[ProductService][DeleteProduct] Producto eliminado con éxito: ID %d", id)
	return nil
}

func RestockProduct(id uint, restockDto dtos.RestockProductDto) error {
	logger.Log.Infof("[ProductService][RestockProduct] Reestockeando producto ID: %d", id)

	if restockDto.PackageCount <= 0 || restockDto.UnitPerPackage <= 0 {
		logger.Log.Warn("[ProductService][RestockProduct] Datos inválidos: PackageCount o UnitPerPackage no pueden ser menores o iguales a 0")
		return errors.New("cantidad de paquetes y unidades por paquete deben ser mayores a 0")
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[ProductService][RestockProduct] Producto no encontrado: ID %d", id)
			return errors.New("producto no encontrado")
		}
		logger.Log.Error("[ProductService][RestockProduct] Error al buscar producto: ", err)
		return err
	}

	// Calcular la cantidad a agregar
	quantityToAdd := restockDto.PackageCount * restockDto.UnitPerPackage
	product.Quantity += quantityToAdd

	// Registrar movimiento de stock (opcional)
	movement := models.StockMovement{
		ProductID:  product.ID,
		Quantity:   quantityToAdd,
		Reason:     restockDto.Reason,
		UnityPrice: restockDto.UnityPrice,
		TotalPrice: restockDto.PackageCount * restockDto.UnityPrice,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&product).Error; err != nil {
			return err
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Log.Error("[ProductService][RestockProduct] Error en transacción: ", err)
		return errors.New("error al crear producto y registrar movimiento inicial")
	}
	logger.Log.Infof("[ProductService][RestockProduct] Producto %s reestockeado. Nueva cantidad: %.2f", product.Name, product.Quantity)

	return nil
}
