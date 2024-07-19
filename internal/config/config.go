package config

import (
	"fmt"
	"os"

	logger "example.com/go-crud/internal/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logDb "gorm.io/gorm/logger"
)

type Config struct {
	DB            *gorm.DB
	ServerAddress string
	DbHost        string
	DbPort        string
	DbUser        string
	DbName        string
	DbPassword    string
}

// getEnv reads an environment variable and returns its value or a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadConfig() *Config {
	logger.InitLogger()
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("No loading for .env file")
	}

	cfg := Config{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbName:     os.Getenv("DB_NAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
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
	// db.AutoMigrate(&models.User{})

	return &Config{
		DB:            db,
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DbHost:        os.Getenv("DB_HOST"),
		DbPort:        os.Getenv("DB_PORT"),
		DbUser:        os.Getenv("DB_USER"),
		DbName:        os.Getenv("DB_NAME"),
		DbPassword:    os.Getenv("DB_PASSWORD"),
	}
}
