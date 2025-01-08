package controllers

import (
	"net/http"
	"peluqueria/internal/dtos"
	"peluqueria/internal/helpers"
	"peluqueria/internal/services"
	"peluqueria/logger"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @Summary Crear producto
// @Description Crea un nuevo producto en el sistema.
// @Tags Productos
// @Accept json
// @Produce json
// @Param product body dtos.CreateProductDto true "Datos del producto"
// @Success 200 {object} dtos.Response{data=nil}
// @Failure 400 {object} dtos.Response{data=nil} "Datos inválidos"
// @Failure 500 {object} dtos.Response{data=nil} "Error interno del servidor"
// @Router /producto [post]
func CreateProduct(c echo.Context) error {
	var product dtos.CreateProductDto
	if err := c.Bind(&product); err != nil {
		logger.Log.Warn("[ProductController][CreateProduct] Error al crear producto: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}
	err := services.CreateProduct(product)
	if err != nil {
		logger.Log.Error("[ProductController][CreateProduct] Error al crear producto: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ProductController][CreateProduct] Producto creado: %s", product.Name)
	return helpers.RespondSuccess(c, "Producto creado", nil)
}

// @Summary Obtener todos los productos
// @Description Devuelve una lista de todos los productos registrados.
// @Tags Productos
// @Produce json
// @Success 200 {object} dtos.Response{data=[]dtos.GetProductDto}
// @Failure 500 {object} dtos.Response{data=nil} "Error interno del servidor"
// @Router /producto [get]
func GetAllProducts(c echo.Context) error {
	products, err := services.GetAllProducts()
	if err != nil {
		logger.Log.Error("[ProductController][GetAllProducts] Error al obtener productos: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.RespondSuccess(c, "Productos obtenidos", products)
}

// @Summary Obtener producto por ID
// @Description Devuelve los datos de un producto específico.
// @Tags Productos
// @Produce json
// @Param id path int true "ID del producto"
// @Success 200 {object} dtos.Response{data=dtos.GetProductDto}
// @Failure 400 {object} dtos.Response{data=nil} "ID inválido"
// @Failure 404 {object} dtos.Response{data=nil} "Producto no encontrado"
// @Router /producto/{id} [get]
func GetProductByID(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ProductController][GetProductByID] Error al obtener producto: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	product, err := services.GetProductByID(uint(productID))
	if err != nil {
		logger.Log.Error("[ProductController][GetProductByID] Error al obtener producto: ", err)
		return helpers.RespondError(c, http.StatusNotFound, err.Error())
	}
	logger.Log.Info("[ProductController][GetProductByID] Producto encontrado: ", product.Name)
	return helpers.RespondSuccess(c, "Producto encontrado", product)
}

// @Summary Actualizar producto
// @Description Actualiza los datos de un producto específico.
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path int true "ID del producto"
// @Param product body dtos.UpdateProductDto true "Datos del producto"
// @Success 200 {object} dtos.Response{data=nil}
// @Failure 400 {object} dtos.Response{data=nil} "Datos o ID inválidos"
// @Failure 500 {object} dtos.Response{data=nil} "Error interno del servidor"
// @Router /producto/{id} [put]
func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ProductController][UpdateProduct] Error al actualizar producto: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var product dtos.UpdateProductDto
	if err := c.Bind(&product); err != nil {
		logger.Log.Warn("[ProductController][UpdateProduct] Error al actualizar producto: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.UpdateProduct(uint(productID), product)
	if err != nil {
		logger.Log.Error("[ProductController][UpdateProduct] Error al actualizar producto: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Infof("[ProductController][UpdateProduct] Producto actualizado: %s", product.Name)
	return helpers.RespondSuccess(c, "Producto actualizado", nil)
}

// @Summary Eliminar producto
// @Description Elimina un producto específico.
// @Tags Productos
// @Produce json
// @Param id path int true "ID del producto"
// @Success 200 {object} dtos.Response{data=nil}
// @Failure 400 {object} dtos.Response{data=nil} "ID inválido"
// @Failure 500 {object} dtos.Response{data=nil} "Error interno del servidor"
// @Router /producto/{id} [delete]
func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.Warn("[ProductController][DeleteProduct] Error al eliminar producto: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	err = services.DeleteProduct(uint(productID))
	if err != nil {
		logger.Log.Error("[ProductController][DeleteProduct] Error al eliminar producto: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Info("[ProductController][DeleteProduct] Producto eliminado: ID ", productID)
	return helpers.RespondSuccess(c, "Producto eliminado", nil)
}

// @Summary Reabastecer producto
// @Description Agrega stock adicional a un producto existente.
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path int true "ID del producto"
// @Param restock body dtos.RestockProductDto true "Datos para el reestock"
// @Success 200 {object} dtos.Response{data=nil}
// @Failure 400 {object} dtos.Response{data=nil} "Datos o ID inválidos"
// @Failure 500 {object} dtos.Response{data=nil} "Error interno del servidor"
// @Router /producto/{id}/restock [post]
func RestockProduct(c echo.Context) error {
	productIDParam := c.Param("id")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	logger.Log.Infof("[ProductController][RestockProduct] Reestock de producto con ID: %d", productID)
	if err != nil {
		logger.Log.Warn("[ProductController][RestockProduct] Error al reestockear producto: ID inválido")
		return helpers.RespondError(c, http.StatusBadRequest, "ID inválido")
	}

	var restockDto dtos.RestockProductDto
	if err := c.Bind(&restockDto); err != nil {
		logger.Log.Warn("[ProductController][RestockProduct] Error al reestockear producto: datos inválidos")
		return helpers.RespondError(c, http.StatusBadRequest, "Datos inválidos")
	}

	err = services.RestockProduct(uint(productID), restockDto)
	if err != nil {
		logger.Log.Error("[ProductController][RestockProduct] Error al reestockear producto: ", err)
		return helpers.RespondError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Log.Info("[ProductController][RestockProduct] Reestock realizado con éxito")
	return helpers.RespondSuccess(c, "Reestock realizado con éxito", nil)
}
