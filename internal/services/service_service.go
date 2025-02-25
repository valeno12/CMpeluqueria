package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"
)

func CreateService(serviceDto dtos.ServiceDto) error {
	logger.Log.Infof("[ServiceService][CreateService] Intentando crear servicio: %s", serviceDto.Name)

	if serviceDto.Name == "" {
		logger.Log.Warn("[ServiceService][CreateService] Nombre requerido")
		return errors.New("el nombre del servicio es requerido")
	}

	if serviceDto.Price <= 0 {
		logger.Log.Warn("[ServiceService][CreateService] Precio inválido")
		return errors.New("el precio del servicio debe ser mayor a 0")
	}

	var existingService models.Service
	if err := database.DB.Where("name = ?", serviceDto.Name).First(&existingService).Error; err == nil {
		logger.Log.Warn("[ServiceService][CreateService] Servicio existente")
		return errors.New("el servicio ya existe")
	}

	service := models.Service{
		Name:                 serviceDto.Name,
		Description:          serviceDto.Description,
		Price:                serviceDto.Price,
		EstimatedTimeMinutes: serviceDto.EstimatedTimeMinutes + serviceDto.EstimatedTimeHours*60,
	}

	if err := database.DB.Create(&service).Error; err != nil {
		logger.Log.Error("[ServiceService][CreateService] Error al crear servicio: ", err)
		return errors.New("error al crear servicio")
	}

	logger.Log.Infof("[ServiceService][CreateService] Servicio creado: %s", serviceDto.Name)
	return nil
}

func GetAllServices() ([]dtos.GetServiceDto, error) {
	logger.Log.Info("[ServiceService][GetAllServices] Obteniendo lista de servicios")

	var services []models.Service
	if err := database.DB.Find(&services).Error; err != nil {
		logger.Log.Error("[ServiceService][GetAllServices] Error al obtener servicios: ", err)
		return nil, errors.New("error al obtener servicios")
	}

	var serviceDtos []dtos.GetServiceDto
	for _, service := range services {
		serviceDtos = append(serviceDtos, dtos.GetServiceDto{
			ID:            service.ID,
			Name:          service.Name,
			Description:   service.Description,
			Price:         service.Price,
			EstimatedTime: service.EstimatedTimeMinutes,
		})
	}

	logger.Log.Infof("[ServiceService][GetAllServices] %d servicios obtenidos", len(serviceDtos))
	return serviceDtos, nil
}

func GetServiceByID(id uint) (dtos.GetServiceDto, error) {
	logger.Log.Infof("[ServiceService][GetServiceByID] Obteniendo servicio con ID: %d", id)

	var service models.Service
	if err := database.DB.Where("id = ?", id).First(&service).Error; err != nil {
		logger.Log.Error("[ServiceService][GetServiceByID] Error al obtener servicio: ", err)
		return dtos.GetServiceDto{}, errors.New("error al obtener servicio")
	}

	serviceDto := dtos.GetServiceDto{
		Name:          service.Name,
		Description:   service.Description,
		Price:         service.Price,
		EstimatedTime: service.EstimatedTimeMinutes,
	}

	logger.Log.Infof("[ServiceService][GetServiceByID] Servicio obtenido: %s", service.Name)
	return serviceDto, nil
}

func UpdateService(id uint, serviceDto dtos.ServiceDto) error {
	logger.Log.Infof("[ServiceService][UpdateService] Actualizando servicio con ID: %d", id)

	var service models.Service
	if err := database.DB.First(&service, id).Error; err != nil {
		logger.Log.Warn("[ServiceService][UpdateService] Servicio no encontrado")
		return errors.New("el servicio no existe")
	}

	if serviceDto.Name != "" {
		service.Name = serviceDto.Name
	}
	if serviceDto.Description != "" {
		service.Description = serviceDto.Description
	}
	if serviceDto.Price > 0 {
		service.Price = serviceDto.Price
	}
	if serviceDto.EstimatedTimeMinutes > 0 || serviceDto.EstimatedTimeHours > 0 {
		service.EstimatedTimeMinutes = serviceDto.EstimatedTimeMinutes + serviceDto.EstimatedTimeHours*60
	}

	if err := database.DB.Save(&service).Error; err != nil {
		logger.Log.Error("[ServiceService][UpdateService] Error al actualizar servicio: ", err)
		return errors.New("error al actualizar servicio")
	}

	logger.Log.Infof("[ServiceService][UpdateService] Servicio actualizado: %s", service.Name)
	return nil
}

func DeleteService(id uint) error {
	logger.Log.Infof("[ServiceService][DeleteService] Eliminando servicio con ID: %d", id)

	// Verificar si el servicio está vinculado a turnos
	var count int64
	if err := database.DB.Model(&models.AppointmentService{}).Where("service_id = ?", id).Count(&count).Error; err != nil {
		logger.Log.Error("[ServiceService][DeleteService] Error al verificar relaciones: ", err)
		return errors.New("error al verificar relaciones del servicio")
	}
	if count > 0 {
		logger.Log.Warn("[ServiceService][DeleteService] Servicio vinculado a turnos, no se puede eliminar")
		return errors.New("el servicio está vinculado a turnos y no se puede eliminar")
	}

	if err := database.DB.Delete(&models.Service{}, id).Error; err != nil {
		logger.Log.Error("[ServiceService][DeleteService] Error al eliminar servicio: ", err)
		return errors.New("error al eliminar servicio")
	}

	logger.Log.Infof("[ServiceService][DeleteService] Servicio eliminado con éxito: ID %d", id)
	return nil
}
