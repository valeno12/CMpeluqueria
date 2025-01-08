package services

import (
	"errors"
	"peluqueria/database"
	"peluqueria/internal/dtos"
	"peluqueria/internal/helpers"
	"peluqueria/internal/models"
	"peluqueria/logger"
)

func GetMonthlyStatistics(month string) (dtos.MonthlyStatisticsDto, error) {
	logger.Log.Infof("[StatisticsService][GetMonthlyStatistics] Generando estadísticas para el mes: %s", month)

	// Parsear el mes para obtener los límites de fechas
	startDate, endDate, err := helpers.ParseMonthFilter(month)
	if err != nil {
		logger.Log.Warn("[StatisticsService][GetMonthlyStatistics] Error al parsear mes: ", err)
		return dtos.MonthlyStatisticsDto{}, err
	}

	var statistics dtos.MonthlyStatisticsDto

	// Calcular ingresos
	var income float64
	if err := database.DB.
		Model(&models.AppointmentService{}).
		Select("COALESCE(SUM(price), 0)").
		Joins("JOIN appointments ON appointments.id = appointment_services.appointment_id").
		Where("appointments.status = 'finalizado' AND appointments.appointment_date BETWEEN ? AND ?", startDate, endDate).
		Scan(&income).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al calcular ingresos: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al calcular ingresos")
	}
	statistics.Incomes = income

	// Calcular egresos (compras de stock)
	var expenses float64
	if err := database.DB.
		Model(&models.StockMovement{}).
		Select("COALESCE(SUM(unity_price * package_count), 0)").
		Where("quantity > 0 AND created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&expenses).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al calcular egresos: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al calcular egresos")
	}
	statistics.Expenses = expenses

	// Agregar las variables para almacenar los totales de cada método de pago
	var debitCount, cashCount, transferCount int64

	// Contar los turnos pagados con débito
	if err := database.DB.Model(&models.Appointment{}).
		Where("payment_method = ? AND status = ? AND appointment_date BETWEEN ? AND ?", "debito", "finalizado", startDate, endDate).
		Count(&debitCount).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al contar turnos con débito: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al contar turnos con débito")
	}

	// Contar los turnos pagados en efectivo
	if err := database.DB.Model(&models.Appointment{}).
		Where("payment_method = ? AND status = ? AND appointment_date BETWEEN ? AND ?", "efectivo", "finalizado", startDate, endDate).
		Count(&cashCount).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al contar turnos con efectivo: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al contar turnos con efectivo")
	}

	// Contar los turnos pagados por transferencia
	if err := database.DB.Model(&models.Appointment{}).
		Where("payment_method = ? AND status = ? AND appointment_date BETWEEN ? AND ?", "transferencia", "finalizado", startDate, endDate).
		Count(&transferCount).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al contar turnos con transferencia: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al contar turnos con transferencia")
	}

	// Asignar los resultados al DTO
	statistics.PaymentMethodBreakdown = dtos.PaymentMethodBreakdownDto{
		Debit:    debitCount,
		Cash:     cashCount,
		Transfer: transferCount,
	}

	// Contar cantidad de turnos realizados
	var appointmentsCount int64
	if err := database.DB.
		Model(&models.Appointment{}).
		Where("status = ? AND appointment_date BETWEEN ? AND ?", "finalizado", startDate, endDate).
		Count(&appointmentsCount).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al contar turnos: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al contar turnos: " + err.Error())
	}
	statistics.AppointmentsCount = int64(appointmentsCount)

	// Contar cantidad de clientes únicos atendidos
	var clientsCount int64
	if err := database.DB.
		Model(&models.Appointment{}).
		Select("COUNT(DISTINCT client_id)").
		Where("status = ? AND appointment_date BETWEEN ? AND ?", "finalizado", startDate, endDate).
		Scan(&clientsCount).Error; err != nil {
		logger.Log.Error("[StatisticsService][GetMonthlyStatistics] Error al contar clientes: ", err)
		return dtos.MonthlyStatisticsDto{}, errors.New("error al contar clientes: " + err.Error())
	}
	statistics.ClientsCount = int64(clientsCount)

	logger.Log.Infof("[StatisticsService][GetMonthlyStatistics] Estadísticas generadas para el mes: %s", month)
	return statistics, nil
}
