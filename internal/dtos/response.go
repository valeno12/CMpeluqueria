package dtos

// Response estructura genérica para respuestas API
type Response struct {
	Status  string      `json:"status" example:"ok"`                             // Ejemplo de estado para éxito o error
	Message string      `json:"message" example:"Operación realizada con éxito"` // Mensaje descriptivo
	Data    interface{} `json:"data,omitempty"`                                  // Datos adicionales
}
