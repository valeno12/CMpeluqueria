package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/models"
	"peluqueria/logger"

	"gorm.io/gorm"
)

func CreateClient(clientDTO dtos.ClientDTO) error {
	logger.Log.Infof("Intentando crear cliente: %s", clientDTO.Name)

	if clientDTO.Name == "" || clientDTO.LastName == "" {
		logger.Log.Warn("Datos faltantes para crear cliente")
		return errors.New("nombre y apellido son obligatorios")
	}

	var existingClient models.Client
	if err := database.DB.Where("name = ? AND last_name = ?", clientDTO.Name, clientDTO.LastName).First(&existingClient).Error; err == nil {
		logger.Log.Warnf("Cliente ya existe: %s %s", clientDTO.Name, clientDTO.LastName)
		return errors.New("cliente ya existe")
	}

	client := models.Client{
		Name:     clientDTO.Name,
		LastName: clientDTO.LastName,
		Phone:    clientDTO.Phone,
		Email:    clientDTO.Email,
	}

	if err := database.DB.Create(&client).Error; err != nil {
		logger.Log.Error("Error al crear cliente: ", err)
		return errors.New("error al crear cliente")
	}

	return nil
}

func GetAllClients() ([]models.Client, error) {
	logger.Log.Info("Intentando obtener clientes")

	var clients []models.Client
	if err := database.DB.Find(&clients).Error; err != nil {
		logger.Log.Error("Error al obtener clientes: ", err)
		return nil, errors.New("error al obtener clientes")
	}

	return clients, nil
}

func GetClientByID(id uint) (models.Client, error) {
	logger.Log.Infof("Intentando obtener cliente con ID: %d", id)

	var client models.Client
	if err := database.DB.Where("id = ?", id).First(&client).Error; err != nil {
		logger.Log.Error("Error al obtener cliente: ", err)
		return models.Client{}, errors.New("error al obtener cliente")
	}

	return client, nil
}

func UpdateClient(id uint, client dtos.ClientDTO) error {
	logger.Log.Infof("Actualizando cliente con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("ID del cliente faltante en actualización")
		return errors.New("el ID del cliente es obligatorio")
	}

	var existingClient models.Client
	if err := database.DB.First(&existingClient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Log.Warnf("Cliente no encontrado para actualizar: ID %d", id)
			return errors.New("el cliente no existe")
		}
		logger.Log.Error("Error al buscar cliente para actualizar: ", err)
		return err
	}

	if client.Name != "" {
		existingClient.Name = client.Name
	}
	if client.LastName != "" {
		existingClient.LastName = client.LastName
	}
	if client.Phone != "" {
		existingClient.Phone = client.Phone
	}
	if client.Email != "" {
		existingClient.Email = client.Email
	}

	if err := database.DB.Save(&existingClient).Error; err != nil {
		logger.Log.Error("Error al actualizar cliente: ", err)
		return errors.New("error al actualizar cliente")
	}

	return nil
}

func DeleteClient(id uint) error {
	logger.Log.Infof("Eliminando cliente con ID: %d", id)
	if id == 0 {
		logger.Log.Warn("ID del usuario faltante en eliminación")
		return errors.New("el ID del usuario es obligatorio")
	}
	if err := database.DB.Delete(&models.Client{}, id).Error; err != nil {
		logger.Log.Error("Error al eliminar cliente: ", err)
		return errors.New("error al eliminar cliente")
	}
	logger.Log.Infof("Cliente eliminado con éxito: ID %d", id)
	return nil
}
