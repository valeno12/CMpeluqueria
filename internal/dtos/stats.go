package dtos

type PaymentMethodBreakdownDto struct {
	Debit    int64 `json:"debit"`    // Cantidad de pagos con débito
	Cash     int64 `json:"cash"`     // Cantidad de pagos en efectivo
	Transfer int64 `json:"transfer"` // Cantidad de pagos por transferencia
}

type MonthlyStatisticsDto struct {
	Incomes                float64                   `json:"income"`
	Expenses               float64                   `json:"expenses"`
	AppointmentsCount      int64                     `json:"appointments_count"`
	ClientsCount           int64                     `json:"clients_count"`
	PaymentMethodBreakdown PaymentMethodBreakdownDto `json:"payment_method_breakdown"`
}
