package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Log es la instancia global de Logrus
var Log = logrus.New()

func InitLogger() {
	// Crear o abrir el archivo de logs
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Error al abrir el archivo de logs: %v", err)
	}

	// Configurar salida múltiple: archivo y consola
	multiWriter := io.MultiWriter(file, os.Stdout)
	Log.SetOutput(multiWriter)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.InfoLevel)

	Log.Info("Logger inicializado con salida múltiple")
}
