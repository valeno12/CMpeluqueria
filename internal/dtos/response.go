package dtos

// Response estructura genérica para respuestas API
type Response struct {
	Status  string      `json:"status" example:"ok"`                             // Ejemplo de estado para éxito o error
	Message string      `json:"message" example:"Operación realizada con éxito"` // Mensaje descriptivo
	Data    interface{} `json:"data,omitempty"`                                  // Datos adicionales
}

// ErrorResponse representa las respuestas de error
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`                  // Siempre "error"
	Message string `json:"message" example:"Descripción del error"` // Detalle del error
}
