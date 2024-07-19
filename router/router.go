package router

import (
	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/handlers"
	"example.com/go-crud/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	userHandler := handlers.NewUserHandler(cfg)

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.GET("", middleware.AuthMiddleware(), userHandler.GetUsers)
		userRouter.POST("", userHandler.CreateUser)
		userRouter.POST("/login", userHandler.LoginUser)
	}

	productRouter := router.Group("/api/v1/products")
	{
		productRouter.GET("", handlers.GetProducts)
	}

	return router
}
