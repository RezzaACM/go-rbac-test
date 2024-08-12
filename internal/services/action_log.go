package services

import (
	"fmt"

	"example.com/go-crud/internal/models"
	logger "example.com/go-crud/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateActionLog(payload models.CreateActionLogRequest, db *gorm.DB) {
	logger.InitLogger()
	actionLog := models.ActionLog{
		UserId:  payload.UserId,
		Action:  payload.Action,
		Details: payload.Details,
	}

	if err := db.Create(&actionLog).Error; err != nil {
		// Handle error
		log.Info().Msg(fmt.Sprintf("Failed to create action log: %v", err))
		panic(err)
	}
}
