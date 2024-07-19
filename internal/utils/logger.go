package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	// Determine the logging level and format based on the environment
	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return
	case "development":
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	default:
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return
	}
}
