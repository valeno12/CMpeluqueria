package database

import (
	"fmt"
	"os"
	"peluqueria/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDatabase() error {
	// Construir el DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Intentos de conexión
	maxRetries := 5
	retryDelay := 10 * time.Second

	for i := 1; i <= maxRetries; i++ {
		logger.Log.Infof("Intentando conectar a la base de datos (Intento %d/%d)...", i, maxRetries)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			logger.Log.Info("Conexión a la base de datos exitosa")
			DB = db
			return nil
		}

		logger.Log.Warnf("Error al conectar con la base de datos: %v", err)
		if i < maxRetries {
			logger.Log.Infof("Reintentando en %v...", retryDelay)
			time.Sleep(retryDelay)
		} else {
			logger.Log.Fatal("No se pudo conectar a la base de datos después de varios intentos")
			return err
		}
	}

	return fmt.Errorf("error inesperado al intentar conectar a la base de datos")
}
