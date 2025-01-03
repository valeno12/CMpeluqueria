package database

import (
	"fmt"
	"os"
	"peluqueria/logger"

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

	// Conectar a la base de datos
	logger.Log.Info("Intentando conectar a la base de datos...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Error al conectar con la base de datos: ", err)
		return err
	}

	logger.Log.Info("Conexión a la base de datos exitosa")
	DB = db
	return nil
}
