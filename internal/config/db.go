package config

import (
	"fmt"

	"example.com/go-crud/internal/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logDb "gorm.io/gorm/logger"
)

var err error

func InitPosgresDB() *gorm.DB {
	//initialize logger
	utils.InitLogger()

	// Initialize the config
	cfg := LoadConfig()

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("No loading for .env file")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbName,
		cfg.DbPassword,
	)

	// Connect to the PostgreSQL database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logDb.Default.LogMode(logDb.Info),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}

	log.Info().Msg("Database connection successfully connected!")
	return db
}
