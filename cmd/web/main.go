package main

import (
	"os"
	"peluqueria/database"
	_ "peluqueria/docs"
	"peluqueria/internal/models"
	"peluqueria/internal/routes"
	"peluqueria/logger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// @title Peluquería API
// @version 1.0
// @description API para la gestión de una peluquería
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Inicializar Logger
	logger.InitLogger()
	logger.Log.Info("Inicializando la aplicación...")

	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error al cargar el archivo .env: ", err)
	}

	// Revisar argumentos
	if len(os.Args) > 1 {
		command := os.Args[1]

		// Ejecutar comando específico
		switch command {
		case "migrate":
			logger.Log.Info("Ejecutando comando: migrate")
			runMigrations()
			return
		case "seed":
			logger.Log.Info("Ejecutando comando: seed")
			runSeeders()
			return
		case "setup":
			logger.Log.Info("Ejecutando comando: setup")
			runSetup()
			return
		default:
			logger.Log.Warn("Comando no reconocido: ", command)
			return
		}
	}

	// Inicializar el servidor en modo normal
	database.InitializeDatabase()
	logger.Log.Info("Base de datos inicializada correctamente")

	e := echo.New()
	routes.RegisterRoutes(e)
	logger.Log.Info("Rutas registradas correctamente")

	// Arrancar el servidor
	port := ":8080"
	logger.Log.Infof("Servidor iniciado en %s", port)
	if err := e.Start(port); err != nil {
		logger.Log.Fatal("Error al iniciar el servidor: ", err)
	}
}

func runMigrations() {
	database.InitializeDatabase()
	logger.Log.Info("Ejecutando migraciones...")
	err := database.DB.AutoMigrate(
		&models.Client{},
		&models.Service{},
		&models.Appointment{},
		&models.AppointmentService{},
		&models.Product{},
		&models.AppointmentProduct{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
	)
	if err != nil {
		logger.Log.Fatal("Error al ejecutar migraciones: ", err)
	}
	logger.Log.Info("Migraciones ejecutadas con éxito")
}

func runSeeders() {
	database.InitializeDatabase()
	logger.Log.Info("Ejecutando seeders...")
	err := database.SeedDatabase(database.DB)
	if err != nil {
		logger.Log.Fatal("Error al ejecutar seeders: ", err)
	}
	logger.Log.Info("Seeders ejecutados con éxito")
}

func runSetup() {
	logger.Log.Info("Ejecutando setup...")
	runMigrations()
	runSeeders()
	logger.Log.Info("Base de datos configurada con éxito (migraciones + seeders)")
}
