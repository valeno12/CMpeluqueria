package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"
)

func CreateService(serviceDto dtos.ServiceDto) error {
	logger.Log.Infof("Intentando crear servicio: %s", serviceDto.Name)
	if serviceDto.Name == "" {
		logger.Log.Warn("[ServiceService][CreateService] Error al crear servicio: nombre requerido")
		return errors.New("el nombre del servicio es requerido")
	}

	if serviceDto.Price <= 0 {
		logger.Log.Warn("[ServiceService][CreateService] Error al crear servicio: precio inválido")
		return errors.New("el precio del servicio debe ser mayor a 0")
	}

	var existingService dtos.ServiceDto
	if err := database.DB.Where("name = ?", serviceDto.Name).First(&existingService).Error; err == nil {
		logger.Log.Warn("[ServiceService][CreateService] Error al crear servicio: servicio existente")
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
	return nil
}

func GetAllServices() ([]models.Service, error) {
	logger.Log.Info("[ServiceService][GetAllServices] Intentando obtener servicios")

	var services []models.Service
	if err := database.DB.Find(&services).Error; err != nil {
		logger.Log.Error("[ServiceService][GetAllServices] Error al obtener servicios: ", err)
		return nil, errors.New("error al obtener servicios")
	}

	return services, nil
}

func GetServiceByID(id uint) (models.Service, error) {
	logger.Log.Infof("[ServiceService][GetServiceByID] Intentando obtener servicio con ID: %d", id)

	var service models.Service
	if err := database.DB.Where("id = ?", id).First(&service).Error; err != nil {
		logger.Log.Error("[ServiceService][GetServiceByID] Error al obtener servicio: ", err)
		return models.Service{}, errors.New("error al obtener servicio")
	}

	return service, nil
}

func UpdateService(id uint, serviceDto dtos.ServiceDto) error {
	logger.Log.Infof("[ServiceService][UpdateService] Actualizando servicio con ID: %d", id)

	service, err := GetServiceByID(id)
	if err != nil {
		return err
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

	return nil
}

func DeleteService(id uint) error {
	logger.Log.Infof("[ServiceService][DeleteService] Eliminando servicio con ID: %d", id)

	service, err := GetServiceByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(&service).Error; err != nil {
		logger.Log.Error("[ServiceService][DeleteService] Error al eliminar servicio: ", err)
		return errors.New("error al eliminar servicio")
	}

	return nil
}
