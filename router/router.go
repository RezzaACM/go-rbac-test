package router

import (
	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/handlers"
	"example.com/go-crud/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Handler
	userHandler := handlers.NewUserHandler(cfg)
	roleHandler := handlers.NewRoleHandler(cfg)
	permissionHandler := handlers.NewPermissionHandler(cfg)

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.GET("", middleware.AuthMiddleware(), userHandler.GetUsers)
		userRouter.POST("", userHandler.CreateUser)
		userRouter.POST("/login", userHandler.LoginUser)
	}

	roleRouter := router.Group("/api/v1/roles")
	{
		roleRouter.GET("", middleware.AuthMiddleware(), roleHandler.GetRoles)
		roleRouter.POST("", middleware.AuthMiddleware(), roleHandler.CreateRole)
		roleRouter.PUT("/:id", middleware.AuthMiddleware(), roleHandler.UpdateRole)
		roleRouter.GET("/:id", middleware.AuthMiddleware(), roleHandler.GetRole)
		roleRouter.DELETE("/:id", middleware.AuthMiddleware(), roleHandler.DeleteRole)
	}

	permissionRouter := router.Group("/api/v1/permissions")
	{
		permissionRouter.GET("", middleware.AuthMiddleware(), permissionHandler.GetPermissions)
		permissionRouter.POST("", middleware.AuthMiddleware(), permissionHandler.CreatePermission)
		permissionRouter.PUT("/:id", middleware.AuthMiddleware(), permissionHandler.UpdatePermission)
		permissionRouter.DELETE("/:id", middleware.AuthMiddleware(), permissionHandler.DeletePermission)
		permissionRouter.GET("/:id", middleware.AuthMiddleware(), permissionHandler.GetPermission)
	}

	productRouter := router.Group("/api/v1/products")
	{
		productRouter.GET("", handlers.GetProducts)
	}

	return router
}
