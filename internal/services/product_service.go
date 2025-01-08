package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"

	"gorm.io/gorm"
)

func CreateProduct(productDto dtos.CreateProductDto) error {
	logger.Log.Infof("[ProductService][CreateProduct] Intentando crear producto: %s", productDto.Name)

	// Validar datos
	if productDto.Name == "" || productDto.Brand == "" || productDto.Unit == "" {
		logger.Log.Warn("[ProductService][CreateProduct] Datos faltantes para crear producto")
		return errors.New("nombre, marca y unidad son obligatorios")
	}
	if productDto.PackageCount <= 0 || productDto.UnitPerPackage <= 0 {
		logger.Log.Warn("[ProductService][CreateProduct] PackageCount o UnitPerPackage inválidos")
		return errors.New("la cantidad de paquetes y unidades por paquete deben ser mayores a 0")
	}

	// Verificar si el producto ya existe
	var existingProduct models.Product
	if err := database.DB.Where("name = ? AND brand = ?", productDto.Name, productDto.Brand).First(&existingProduct).Error; err == nil {
		logger.Log.Warnf("[ProductService][CreateProduct] Producto ya existente: %s %s", productDto.Name, productDto.Brand)
		return errors.New("el producto ya existe")
	}

	// Crear el producto
	product := models.Product{
		Name:          productDto.Name,
		Unit:          productDto.Unit,
		Brand:         productDto.Brand,
		Quantity:      productDto.PackageCount * productDto.UnitPerPackage,
		LowStockAlert: productDto.LowStockAlert,
	}

	movement := models.StockMovement{
		ProductID:      product.ID,
		ProductUnit:    productDto.Unit,
		PackageCount:   &productDto.PackageCount,
		UnitPerPackage: &productDto.UnitPerPackage,
		Quantity:       productDto.PackageCount * productDto.UnitPerPackage,
		Reason:         "Inventario inicial",
		UnityPrice:     &productDto.UnityPrice,
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
		logger.Log.Error("[ProductService][CreateProduct] Error al crear producto: ", err)
		return errors.New("error al crear el producto")
	}

	logger.Log.Infof("[ProductService][CreateProduct] Producto creado con éxito: %s", product.Name)
	return nil
}

func GetAllProducts() ([]dtos.GetProductDto, error) {
	logger.Log.Info("[ProductService][GetAllProducts] Obteniendo todos los productos")

	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		logger.Log.Error("[ProductService][GetAllProducts] Error al obtener productos: ", err)
		return nil, errors.New("error al obtener productos")
	}
	logger.Log.Infof("[ProductService][GetAllProducts] %d productos encontrados", len(products))
	var productDto []dtos.GetProductDto
	for _, product := range products {
		productDto = append(productDto, dtos.GetProductDto{
			ID:       product.ID,
			Name:     product.Name,
			Brand:    product.Brand,
			Unit:     product.Unit,
			Quantity: product.Quantity,
		})
	}
	return productDto, nil
}

func GetProductByID(id uint) (dtos.GetProductDto, error) {
	logger.Log.Infof("[ProductService][GetProductByID] Intentando obtener producto con ID: %d", id)

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[ProductService][GetProductByID] Producto no encontrado: ID %d", id)
			return dtos.GetProductDto{}, errors.New("producto no encontrado")
		}
		logger.Log.Error("[ProductService][GetProductByID] Error al obtener producto: ", err)
		return dtos.GetProductDto{}, err
	}

	logger.Log.Infof("[ProductService][GetProductByID] Producto encontrado: %s", product.Name)
	productDto := dtos.GetProductDto{
		ID:       product.ID,
		Name:     product.Name,
		Brand:    product.Brand,
		Unit:     product.Unit,
		Quantity: product.Quantity,
	}
	return productDto, nil
}

func UpdateProduct(id uint, productDto dtos.UpdateProductDto) error {
	logger.Log.Infof("[ProductService][UpdateProduct] Actualizando producto con ID: %d", id)

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[ProductService][UpdateProduct] Producto no encontrado: ID %d", id)
			return errors.New("producto no encontrado")
		}
		logger.Log.Error("[ProductService][UpdateProduct] Error al buscar producto: ", err)
		return err
	}

	// Actualizar los campos necesarios
	if productDto.Name != "" {
		product.Name = productDto.Name
	}
	if productDto.Brand != "" {
		product.Brand = productDto.Brand
	}
	if productDto.Unit != "" {
		product.Unit = productDto.Unit
	}
	if productDto.LowStockAlert > 0 {
		product.LowStockAlert = productDto.LowStockAlert
	}

	if err := database.DB.Save(&product).Error; err != nil {
		logger.Log.Error("[ProductService][UpdateProduct] Error al actualizar producto: ", err)
		return errors.New("error al actualizar producto")
	}

	logger.Log.Infof("[ProductService][UpdateProduct] Producto actualizado con éxito: %s", product.Name)
	return nil
}

func DeleteProduct(id uint) error {
	logger.Log.Infof("[ProductService][DeleteProduct] Eliminando producto con ID: %d", id)

	// Verificar si el producto está asociado a un turno
	var count int64
	if err := database.DB.Model(&models.AppointmentProduct{}).Where("product_id = ?", id).Count(&count).Error; err != nil {
		logger.Log.Error("[ProductService][DeleteProduct] Error al verificar asociaciones: ", err)
		return errors.New("error al verificar asociaciones")
	}
	if count > 0 {
		logger.Log.Warn("[ProductService][DeleteProduct] El producto está asociado a turnos y no se puede eliminar")
		return errors.New("el producto está asociado a turnos")
	}

	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		logger.Log.Error("[ProductService][DeleteProduct] Error al eliminar producto: ", err)
		return errors.New("error al eliminar producto")
	}

	logger.Log.Infof("[ProductService][DeleteProduct] Producto eliminado con éxito: ID %d", id)
	return nil
}

func RestockProduct(id uint, restockDto dtos.RestockProductDto) error {
	logger.Log.Infof("[ProductService][RestockProduct] Reestockeando producto con ID: %d", id)

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

	// Registrar movimiento de stock
	movement := models.StockMovement{
		ProductID:      product.ID,
		ProductUnit:    product.Unit,
		PackageCount:   &restockDto.PackageCount,
		UnitPerPackage: &restockDto.UnitPerPackage,
		Quantity:       quantityToAdd,
		Reason:         restockDto.Reason,
		UnityPrice:     &restockDto.UnityPrice,
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
		logger.Log.Error("[ProductService][RestockProduct] Error en la transacción: ", err)
		return errors.New("error al reestockear producto")
	}

	logger.Log.Infof("[ProductService][RestockProduct] Producto reestockeado con éxito: %s", product.Name)
	return nil
}
