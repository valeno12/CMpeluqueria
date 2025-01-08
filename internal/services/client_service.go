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
	logger.Log.Infof("[ClientService][CreateClient] Intentando crear cliente: %s %s", clientDTO.Name, clientDTO.LastName)

	if clientDTO.Name == "" || clientDTO.LastName == "" {
		logger.Log.Warn("[ClientService][CreateClient] Datos faltantes para crear cliente")
		return errors.New("nombre y apellido son obligatorios")
	}

	var existingClient models.Client
	if err := database.DB.Where("name = ? AND last_name = ?", clientDTO.Name, clientDTO.LastName).First(&existingClient).Error; err == nil {
		logger.Log.Warnf("[ClientService][CreateClient] Cliente ya existe: %s %s", clientDTO.Name, clientDTO.LastName)
		return errors.New("cliente ya existe")
	}

	client := models.Client{
		Name:     clientDTO.Name,
		LastName: clientDTO.LastName,
		Phone:    clientDTO.Phone,
		Email:    clientDTO.Email,
	}

	if err := database.DB.Create(&client).Error; err != nil {
		logger.Log.Error("[ClientService][CreateClient] Error al crear cliente: ", err)
		return errors.New("error al crear cliente")
	}

	logger.Log.Infof("[ClientService][CreateClient] Cliente creado con éxito: %s %s", clientDTO.Name, clientDTO.LastName)
	return nil
}

func GetAllClients() ([]dtos.GetClientDto, error) {
	logger.Log.Info("[ClientService][GetAllClients] Intentando obtener clientes")

	var clients []models.Client
	if err := database.DB.Find(&clients).Error; err != nil {
		logger.Log.Error("[ClientService][GetAllClients] Error al obtener clientes: ", err)
		return nil, errors.New("error al obtener clientes")
	}

	var clientsDto []dtos.GetClientDto
	for _, client := range clients {
		clientsDto = append(clientsDto, dtos.GetClientDto{
			ID:       client.ID,
			Name:     client.Name,
			LastName: client.LastName,
			Phone:    client.Phone,
			Email:    client.Email,
		})
	}

	logger.Log.Infof("[ClientService][GetAllClients] Clientes obtenidos: %d", len(clients))
	return clientsDto, nil
}

func GetClientByID(id uint) (dtos.GetClientDto, error) {
	logger.Log.Infof("[ClientService][GetClientByID] Intentando obtener cliente con ID: %d", id)

	var client models.Client

	err := database.DB.
		Where("id = ?", id).
		First(&client).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[ClientService][GetClientByID] Cliente no encontrado: ID %d", id)
			return dtos.GetClientDto{}, errors.New("cliente no encontrado")
		}
		logger.Log.Error("[ClientService][GetClientByID] Error al obtener cliente: ", err)
		return dtos.GetClientDto{}, errors.New("error al obtener cliente")
	}

	var appointments []models.Appointment

	if err := database.DB.
		Where("client_id = ?", id).
		Find(&appointments).Error; err != nil {
		logger.Log.Error("[ClientService][GetClientByID] Error al obtener turnos del cliente: ", err)
		return dtos.GetClientDto{}, errors.New("error al obtener turnos del cliente")
	}

	var appointmentDtos []dtos.ClientAppointmentDto

	for _, appointment := range appointments {
		appointmentDtos = append(appointmentDtos, dtos.ClientAppointmentDto{
			ID:              appointment.ID,
			AppointmentDate: appointment.AppointmentDate.Format("02/01/2006 15:04"),
			Status:          appointment.Status,
		})
	}
	clientDto := dtos.GetClientDto{
		ID:           client.ID,
		Name:         client.Name,
		LastName:     client.LastName,
		Phone:        client.Phone,
		Email:        client.Email,
		Appointments: appointmentDtos,
	}

	logger.Log.Infof("[ClientService][GetClientByID] Cliente obtenido con éxito: ID %d", id)
	return clientDto, nil
}

func UpdateClient(id uint, clientDTO dtos.ClientDTO) error {
	logger.Log.Infof("[ClientService][UpdateClient] Actualizando cliente con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("[ClientService][UpdateClient] ID del cliente faltante en actualización")
		return errors.New("el ID del cliente es obligatorio")
	}

	var existingClient models.Client
	if err := database.DB.First(&existingClient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("[ClientService][UpdateClient] Cliente no encontrado para actualizar: ID %d", id)
			return errors.New("el cliente no existe")
		}
		logger.Log.Error("[ClientService][UpdateClient] Error al buscar cliente para actualizar: ", err)
		return err
	}

	if clientDTO.Name != "" {
		existingClient.Name = clientDTO.Name
	}
	if clientDTO.LastName != "" {
		existingClient.LastName = clientDTO.LastName
	}
	if clientDTO.Phone != "" {
		existingClient.Phone = clientDTO.Phone
	}
	if clientDTO.Email != "" {
		existingClient.Email = clientDTO.Email
	}

	if err := database.DB.Save(&existingClient).Error; err != nil {
		logger.Log.Error("[ClientService][UpdateClient] Error al actualizar cliente: ", err)
		return errors.New("error al actualizar cliente")
	}

	logger.Log.Infof("[ClientService][UpdateClient] Cliente actualizado con éxito: ID %d", id)
	return nil
}

func DeleteClient(id uint) error {
	logger.Log.Infof("[ClientService][DeleteClient] Intentando eliminar cliente con ID: %d", id)

	if id == 0 {
		logger.Log.Warn("[ClientService][DeleteClient] ID del cliente faltante en eliminación")
		return errors.New("el ID del cliente es obligatorio")
	}

	// Verificar si el cliente tiene citas asignadas
	var appointmentCount int64
	if err := database.DB.Model(&models.Appointment{}).Where("client_id = ?", id).Count(&appointmentCount).Error; err != nil {
		logger.Log.Error("[ClientService][DeleteClient] Error al verificar citas del cliente: ", err)
		return errors.New("error al verificar citas del cliente")
	}

	if appointmentCount > 0 {
		logger.Log.Warnf("[ClientService][DeleteClient] El cliente tiene %d citas asignadas y no puede ser eliminado", appointmentCount)
		return errors.New("no se puede eliminar un cliente con citas asignadas")
	}

	if err := database.DB.Delete(&models.Client{}, id).Error; err != nil {
		logger.Log.Error("[ClientService][DeleteClient] Error al eliminar cliente: ", err)
		return errors.New("error al eliminar cliente")
	}

	logger.Log.Infof("[ClientService][DeleteClient] Cliente eliminado con éxito: ID %d", id)
	return nil
}
