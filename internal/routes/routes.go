package routes

import (
	"peluqueria/internal/controllers"
	"peluqueria/middlewares"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const prefix = "/api/v1"

func RegisterRoutes(e *echo.Echo) {
	e.Use(middlewares.RouteLogger)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Rutas públicas
	e.POST(prefix+"/login", controllers.Login) // Iniciar sesión

	userGroup := e.Group(prefix+"/usuarios", middlewares.JWTMiddleware)
	userGroup.POST("", controllers.CreateUser, middlewares.PermissionMiddleware("create_user"))
	userGroup.GET("", controllers.GetAllUsers)
	userGroup.GET("/:id", controllers.GetUserByID)
	userGroup.PUT("/:id", controllers.UpdateUser)
	userGroup.DELETE("/:id", controllers.DeleteUser, middlewares.PermissionMiddleware("delete_user"))

	roleGroup := e.Group(prefix+"/rol", middlewares.JWTMiddleware)
	roleGroup.POST("", controllers.CreateRole, middlewares.PermissionMiddleware("create_role"))
	roleGroup.GET("", controllers.GetAllRoles)
	roleGroup.GET("/:id", controllers.GetRoleByID)
	roleGroup.PUT("/:id", controllers.UpdateRole, middlewares.PermissionMiddleware("update_role"))
	roleGroup.DELETE("/:id", controllers.DeleteRole, middlewares.PermissionMiddleware("delete_role"))

	clientGroup := e.Group(prefix+"/cliente", middlewares.JWTMiddleware)
	clientGroup.POST("", controllers.CreateClient, middlewares.PermissionMiddleware("create_client"))
	clientGroup.GET("", controllers.GetAllClients)
	clientGroup.GET("/:id", controllers.GetClientByID)
	clientGroup.PUT("/:id", controllers.UpdateClient, middlewares.PermissionMiddleware("update_client"))
	clientGroup.DELETE("/:id", controllers.DeleteClient, middlewares.PermissionMiddleware("delete_client"))

	productGroup := e.Group(prefix+"/producto", middlewares.JWTMiddleware)
	productGroup.POST("", controllers.CreateProduct, middlewares.PermissionMiddleware("create_product"))
	productGroup.GET("", controllers.GetAllProducts)
	productGroup.GET("/:id", controllers.GetProductByID)
	productGroup.PUT("/:id", controllers.UpdateProduct, middlewares.PermissionMiddleware("update_product"))
	productGroup.DELETE("/:id", controllers.DeleteProduct, middlewares.PermissionMiddleware("delete_product"))
	productGroup.POST("/:id/restock", controllers.RestockProduct, middlewares.PermissionMiddleware("restock_product"))

	stockGroup := e.Group(prefix+"/stock-movements", middlewares.JWTMiddleware)
	stockGroup.GET("", controllers.GetStockMovements)                      // Todos los movimientos
	stockGroup.GET("/product/:id", controllers.GetStockMovementsByProduct) // Movimientos por producto

	serviceGroup := e.Group(prefix+"/servicio", middlewares.JWTMiddleware)
	serviceGroup.POST("", controllers.CreateService, middlewares.PermissionMiddleware("create_service"))
	serviceGroup.GET("", controllers.GetAllServices)
	serviceGroup.GET("/:id", controllers.GetServiceByID)
	serviceGroup.PUT("/:id", controllers.UpdateService, middlewares.PermissionMiddleware("update_service"))
	serviceGroup.DELETE("/:id", controllers.DeleteService, middlewares.PermissionMiddleware("delete_service"))

	appointmentGroup := e.Group(prefix+"/turno", middlewares.JWTMiddleware)
	appointmentGroup.POST("", controllers.CreateAppointment, middlewares.PermissionMiddleware("create_appointment"))
	appointmentGroup.GET("", controllers.GetAllAppointments)
	appointmentGroup.GET("/:id", controllers.GetAppointmentByID)
	appointmentGroup.PUT("/:id", controllers.UpdateAppointment, middlewares.PermissionMiddleware("update_appointment"))
	appointmentGroup.PUT("/:id/products", controllers.UpdateAppointmentProducts, middlewares.PermissionMiddleware("update_appointment"))
	appointmentGroup.DELETE("/:id", controllers.DeleteAppointment, middlewares.PermissionMiddleware("delete_appointment"))
	appointmentGroup.PUT("/:id/finalizar", controllers.FinalizeAppointment)
}
