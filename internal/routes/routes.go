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
}
