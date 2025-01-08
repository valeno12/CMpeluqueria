package services

import (
	"errors"
	"fmt"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/helpers"
	"peluqueria/internal/models"
	"peluqueria/logger"

	"gorm.io/gorm"
)

func CreateAppointment(appointmentDto dtos.CreateAppointmentDto) error {
	logger.Log.Infof("[AppointmentService][CreateAppointment] Creando cita para cliente ID: %d", appointmentDto.ClientID)

	appointmentDate, err := helpers.ParseCustomDate(appointmentDto.AppointmentDate)
	if err != nil {
		logger.Log.Warn("[AppointmentService][CreateAppointment] Error al parsear fecha: ", err)
		return err
	}

	// Validar que el cliente existe
	var client models.Client
	if err := database.DB.First(&client, appointmentDto.ClientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[AppointmentService][CreateAppointment] Cliente no encontrado: ID %d", appointmentDto.ClientID)
			return errors.New("cliente no encontrado")
		}
		logger.Log.Error("[AppointmentService][CreateAppointment] Error al buscar cliente: ", err)
		return errors.New("error al buscar cliente")
	}

	// Validar que los servicios existen
	var services []models.Service
	if len(appointmentDto.ServiceIds) > 0 {
		if err := database.DB.Where("id IN ?", appointmentDto.ServiceIds).Find(&services).Error; err != nil {
			logger.Log.Error("[AppointmentService][CreateAppointment] Error al buscar servicios: ", err)
			return errors.New("error al buscar servicios")
		}
		if len(services) != len(appointmentDto.ServiceIds) {
			logger.Log.Warn("[AppointmentService][CreateAppointment] Uno o más servicios no existen")
			return errors.New("uno o más servicios no existen")
		}
	}

	// Crear la cita
	appointment := models.Appointment{
		ClientID:        appointmentDto.ClientID,
		AppointmentDate: appointmentDate,
		Status:          "pendiente", // Estado inicial
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&appointment).Error; err != nil {
			return err
		}

		// Asociar servicios al appointment
		for _, service := range services {
			appointmentService := models.AppointmentService{
				AppointmentID: appointment.ID,
				ServiceID:     service.ID,
				Price:         service.Price,
			}
			if err := tx.Create(&appointmentService).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		logger.Log.Error("[AppointmentService][CreateAppointment] Error al crear cita: ", err)
		return errors.New("error al crear cita")
	}

	logger.Log.Infof("[AppointmentService][CreateAppointment] Cita creada con éxito: ID %d", appointment.ID)
	return nil
}

func GetAllAppointments(clientID, status, startDate, endDate string) ([]dtos.AllAppointmentDto, error) {
	logger.Log.Info("[AppointmentService][GetAllAppointments] Obteniendo todas las citas con filtros")

	var appointments []models.Appointment
	query := database.DB.Preload("Client").Preload("AppointmentServices.Service")

	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("appointment_date BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&appointments).Error; err != nil {
		logger.Log.Error("[AppointmentService][GetAllAppointments] Error al obtener turnos: ", err)
		return nil, errors.New("error al obtener citas")
	}

	var appointmentsDto []dtos.AllAppointmentDto
	for _, appointment := range appointments {
		var timeInMinutes uint
		for _, appService := range appointment.AppointmentServices {
			timeInMinutes += appService.Service.EstimatedTimeMinutes
		}
		appointmentsDto = append(appointmentsDto, dtos.AllAppointmentDto{
			ID:                   appointment.ID,
			ClientID:             appointment.ClientID,
			ClientName:           fmt.Sprintf("%s %s", appointment.Client.Name, appointment.Client.LastName),
			Status:               appointment.Status,
			AppointmentDate:      appointment.AppointmentDate.Format("02/01/2006 15:04"),
			EstimatedTimeMinutes: timeInMinutes,
		})
	}

	return appointmentsDto, nil
}

func GetAppointmentByID(id uint) (dtos.AppointmentByIDDto, error) {
	logger.Log.Infof("[AppointmentService][GetAppointmentByID] Intentando obtener turno con ID: %d", id)
	var appointment models.Appointment
	err := database.DB.
		Preload("Client").
		Preload("AppointmentServices.Service").
		Preload("AppointmentProducts.Product").
		First(&appointment, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[AppointmentService][GetAppointmentByID] Turno no encontrado: ID %d", id)
			return dtos.AppointmentByIDDto{}, errors.New("turno no encontrado")
		}
		logger.Log.Error("[AppointmentService][GetAppointmentByID] Error al obtener turno: ", err)
		return dtos.AppointmentByIDDto{}, errors.New("error al obtener turno")
	}

	var services []dtos.AppointmentServiceDto
	for _, appService := range appointment.AppointmentServices {
		services = append(services, dtos.AppointmentServiceDto{
			ServiceID:            appService.Service.ID,
			ServiceName:          appService.Service.Name,
			Price:                appService.Price,
			EstimatedTimeMinutes: appService.Service.EstimatedTimeMinutes,
		})
	}

	var products []dtos.AppointmentProductDto
	for _, appProduct := range appointment.AppointmentProducts {
		products = append(products, dtos.AppointmentProductDto{
			ProductID: appProduct.Product.ID,
			Name:      appProduct.Product.Name,
			Quantity:  appProduct.Quantity,
			Unit:      appProduct.Product.Unit,
		})
	}

	appointmentDto := dtos.AppointmentByIDDto{
		ID:              appointment.ID,
		ClientID:        appointment.ClientID,
		ClientName:      fmt.Sprintf("%s %s", appointment.Client.Name, appointment.Client.LastName),
		Status:          appointment.Status,
		AppointmentDate: appointment.AppointmentDate.Format("02/01/2006 15:04"),
		Services:        services,
		Products:        products,
		CreatedAt:       appointment.CreatedAt,
		UpdatedAt:       appointment.UpdatedAt,
	}

	logger.Log.Infof("[AppointmentService][GetAppointmentByID] Turno obtenido con éxito: ID %d", id)
	return appointmentDto, nil
}

func UpdateAppointment(id uint, appointmentDto dtos.CreateAppointmentDto) error {
	logger.Log.Infof("[AppointmentService][UpdateAppointment] Actualizando turno con ID: %d", id)
	if id == 0 {
		logger.Log.Warn("[AppointmentService][UpdateAppointment] ID del turno faltante en actualización")
		return errors.New("el ID del turno es obligatorio")
	}

	// Iniciar transacción
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var existingAppointment models.Appointment
		if err := tx.First(&existingAppointment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Log.Warnf("[AppointmentService][UpdateAppointment] Turno no encontrado para actualizar: ID %d", id)
				return errors.New("el Turno no existe")
			}
			logger.Log.Error("[AppointmentService][UpdateAppointment] Error al buscar Turno para actualizar: ", err)
			return err
		}
		if appointmentDto.AppointmentDate != "" {
			appointmentDate, err := helpers.ParseCustomDate(appointmentDto.AppointmentDate)
			if err != nil {
				logger.Log.Warn("[AppointmentService][UpdateAppointment] Error al parsear fecha: ", err)
				return err
			}
			existingAppointment.AppointmentDate = appointmentDate
		}

		if appointmentDto.ClientID != 0 {
			if err := database.DB.Select("id").First(&models.Client{}, appointmentDto.ClientID).Error; err != nil {
				logger.Log.Warn("[AppointmentService][UpdateAppointment] cliente no encontrado: ")
				return errors.New("cliente no encontrado")
			}
			existingAppointment.ClientID = appointmentDto.ClientID
		}

		if err := tx.Save(&existingAppointment).Error; err != nil {
			logger.Log.Error("[AppointmentService][UpdateAppointment] Error al actualizar fecha: ", err)
			return errors.New("error al actualizar la fecha del turno")
		}

		if len(appointmentDto.ServiceIds) > 0 {
			logger.Log.Infof("[AppointmentService][UpdateAppointment] Actualizando servicios del turno")
			if err := tx.Unscoped().Where("appointment_id = ?", existingAppointment.ID).Delete(&models.AppointmentService{}).Error; err != nil {
				logger.Log.Error("[AppointmentService][UpdateAppointment] Error al eliminar servicios antiguos: ", err)
				return errors.New("error al eliminar servicios antiguos")
			}

			for _, serviceId := range appointmentDto.ServiceIds {
				var service models.Service
				if err := tx.First(&service, serviceId).Error; err != nil {
					logger.Log.Warn("[AppointmentService][UpdateAppointment] Servicio no encontrado: ID ", serviceId)
					return errors.New("servicio no encontrado")
				}

				appointmentService := models.AppointmentService{
					AppointmentID: existingAppointment.ID,
					ServiceID:     serviceId,
					Price:         service.Price,
				}

				if err := tx.Create(&appointmentService).Error; err != nil {
					logger.Log.Error("[AppointmentService][UpdateAppointment] Error al asignar servicios al turno: ", err)
					return errors.New("error al asignar servicios al turno")
				}
			}
		}
		return nil
	})
	if err != nil {
		logger.Log.Error("[AppointmentService][UpdateAppointment] Error en transacción: ", err)
		return errors.New("error al actualizar el turno: " + err.Error())
	}
	logger.Log.Infof("[AppointmentService][UpdateAppointment] Turno actualizado con éxito")
	return nil
}

func DeleteAppointment(id uint) error {
	logger.Log.Infof("eliminando turno con id: %v", id)
	if id == 0 {
		logger.Log.Warn("ID del turno faltante en eliminación")
		return errors.New("el ID del turno es obligatorio")
	}
	if err := database.DB.Delete(&models.Appointment{}, id).Error; err != nil {
		logger.Log.Error("Error al eliminar turno: ", err)
		return errors.New("error al eliminar turno")
	}
	if err := database.DB.Where("appointment_id = ?", id).Delete(&models.AppointmentService{}).Error; err != nil {
		logger.Log.Error("Error al servicios asociados al turno: ", err)
		return errors.New("error al eliminar servicios asociados al turno")
	}
	if err := database.DB.Where("appointment_id = ?", id).Delete(&models.AppointmentProduct{}).Error; err != nil {
		logger.Log.Error("Error al servicios productos al turno: ", err)
		return errors.New("error al eliminar servicios productos al turno")
	}
	logger.Log.Infof("turno eliminado con éxito: ID %d", id)
	return nil
}

func FinalizeAppointment(id uint, finalizeDto dtos.FinalizeAppointmentDto) error {
	logger.Log.Infof("[AppointmentService][FinalizeAppointment] Finalizando turno con ID: %d", id)

	if finalizeDto.PaymentMethod == "" {
		logger.Log.Warn("[AppointmentService][FinalizeAppointment] Método de pago faltante")
		return errors.New("el método de pago es obligatorio")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var appointment models.Appointment
		if err := tx.First(&appointment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Log.Warnf("[AppointmentService][FinalizeAppointment] Turno no encontrado: ID %d", id)
				return errors.New("turno no encontrado")
			}
			logger.Log.Error("[AppointmentService][FinalizeAppointment] Error al buscar turno: ", err)
			return err
		}

		if appointment.Status == "finalizado" {
			logger.Log.Warn("[AppointmentService][FinalizeAppointment] El turno ya está finalizado")
			return errors.New("el turno ya está finalizado")
		}

		// Actualizar el estado del turno
		appointment.Status = "finalizado"
		appointment.PaymentMethod = finalizeDto.PaymentMethod
		if err := tx.Save(&appointment).Error; err != nil {
			logger.Log.Error("[AppointmentService][FinalizeAppointment] Error al actualizar estado del turno: ", err)
			return errors.New("error al actualizar estado del turno")
		}

		// Registrar productos utilizados (si se incluyen)
		if len(finalizeDto.Products) > 0 {
			for _, productDto := range finalizeDto.Products {
				var product models.Product
				if err := tx.First(&product, productDto.ProductID).Error; err != nil {
					logger.Log.Warnf("[AppointmentService][FinalizeAppointment] Producto no encontrado: ID %d", productDto.ProductID)
					return errors.New("producto no encontrado")
				}

				if product.Quantity < productDto.Quantity {
					logger.Log.Warnf("[AppointmentService][FinalizeAppointment] Stock insuficiente para producto: ID %d", productDto.ProductID)
					return errors.New("stock insuficiente para producto")
				}

				// Actualizar el stock del producto
				product.Quantity -= productDto.Quantity
				if err := tx.Save(&product).Error; err != nil {
					logger.Log.Error("[AppointmentService][FinalizeAppointment] Error al actualizar stock del producto: ", err)
					return errors.New("error al actualizar stock del producto")
				}

				// Registrar el producto en AppointmentProducts
				appointmentProduct := models.AppointmentProduct{
					AppointmentID: appointment.ID,
					ProductID:     product.ID,
					Quantity:      productDto.Quantity,
				}
				if err := tx.Create(&appointmentProduct).Error; err != nil {
					logger.Log.Error("[AppointmentService][FinalizeAppointment] Error al registrar producto en turno: ", err)
					return errors.New("error al registrar producto en turno")
				}

				// Registrar el movimiento de stock
				stockMovement := models.StockMovement{
					ProductID:   product.ID,
					Quantity:    -productDto.Quantity, // Negativo para salida
					ProductUnit: product.Unit,
					Reason:      fmt.Sprintf("Utilización en turno ID %d", appointment.ID),
				}
				if err := tx.Create(&stockMovement).Error; err != nil {
					logger.Log.Error("[AppointmentService][FinalizeAppointment] Error al registrar movimiento de stock: ", err)
					return errors.New("error al registrar movimiento de stock")
				}
			}
		}

		logger.Log.Infof("[AppointmentService][FinalizeAppointment] Turno finalizado con éxito: ID %d", id)
		return nil
	})
}

func UpdateAppointmentProducts(appointmentID uint, dto dtos.UpdateAppointmentProductsDto) error {
	logger.Log.Infof("[AppointmentService][UpdateAppointmentProducts] Actualizando productos para turno ID: %d", appointmentID)

	var appointment models.Appointment
	if err := database.DB.Preload("AppointmentProducts").First(&appointment, appointmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[AppointmentService][UpdateAppointmentProducts] Turno no encontrado: ID %d", appointmentID)
			return errors.New("turno no encontrado")
		}
		logger.Log.Error("[AppointmentService][UpdateAppointmentProducts] Error al buscar turno: ", err)
		return errors.New("error al buscar turno")
	}

	if appointment.Status != "finalizado" {
		logger.Log.Warnf("[AppointmentService][UpdateAppointmentProducts] El turno no está finalizado: ID %d", appointmentID)
		return errors.New("solo se pueden actualizar productos de un turno finalizado")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Eliminar productos existentes
		if err := tx.Where("appointment_id = ?", appointmentID).Delete(&models.AppointmentProduct{}).Error; err != nil {
			logger.Log.Error("[AppointmentService][UpdateAppointmentProducts] Error al eliminar productos existentes: ", err)
			return errors.New("error al eliminar productos existentes")
		}

		// Agregar nuevos productos
		for _, product := range dto.Products {
			var dbProduct models.Product
			if err := tx.First(&dbProduct, product.ProductID).Error; err != nil {
				logger.Log.Warnf("[AppointmentService][UpdateAppointmentProducts] Producto no encontrado: ID %d", product.ProductID)
				return errors.New("producto no encontrado")
			}

			// Crear relación de producto
			appointmentProduct := models.AppointmentProduct{
				AppointmentID: appointmentID,
				ProductID:     product.ProductID,
				Quantity:      product.Quantity,
			}

			if err := tx.Create(&appointmentProduct).Error; err != nil {
				logger.Log.Error("[AppointmentService][UpdateAppointmentProducts] Error al crear producto: ", err)
				return errors.New("error al actualizar productos del turno")
			}
		}
		return nil
	})

	if err != nil {
		logger.Log.Error("[AppointmentService][UpdateAppointmentProducts] Error en transacción: ", err)
		return err
	}

	logger.Log.Infof("[AppointmentService][UpdateAppointmentProducts] Productos actualizados con éxito para turno ID: %d", appointmentID)
	return nil
}
