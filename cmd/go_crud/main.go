package main

import (
	"os"
	"os/signal"
	"syscall"

	config "example.com/go-crud/internal/config"
	"example.com/go-crud/internal/server"
	logger "example.com/go-crud/internal/utils"
	"example.com/go-crud/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	// initialize the logger
	logger.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	router := router.SetupRouter(cfg)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello world!"})
	})

	// Create a new HTTP Server
	srv := server.NewServer(cfg, router)

	// Channel to listen for OS singals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal().Err(err).Msgf("Could not listen on%s", os.Getenv("SERVER_ADDRESS"))
		}
	}()
	log.Info().Msgf("Server is ready to handle requests at http://localhost%s", os.Getenv("SERVER_ADDRESS"))

	// Block until we receive a signal
	<-quit

	// Attempt a graceful shutdown
	if err := srv.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
