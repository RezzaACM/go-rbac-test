package handlers

import (
	"net/http"
	"strconv"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActionLogHandler struct {
	db *gorm.DB
}

func NewActionLogHandler(db *config.Config) *ActionLogHandler {
	return &ActionLogHandler{db: db.DB}
}

func (h *ActionLogHandler) GetActionLogs(c *gin.Context) {
	var actionLogs []models.ActionLog
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	query := h.db.Preload("User")
	if err := query.Limit(limit).Offset(offset).Find(&actionLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.ActionLog{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}
	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, actionLogs, meta)
}

func (h *ActionLogHandler) GetActionLog(c *gin.Context) {
	var actionLog models.ActionLog
	if err := h.db.Preload("User").First(&actionLog, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, actionLog, nil)
}

func (h *ActionLogHandler) CreateActionLog(c *gin.Context) {
	var actionLog models.CreateActionLogRequest
	if err := c.ShouldBindJSON(&actionLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Create(&actionLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, actionLog)
}

func (h *ActionLogHandler) UpdateActionLog(c *gin.Context) {
	var actionLog models.UpdateActionLogRequest
	var actionLogModel models.ActionLog

	if err := c.ShouldBindJSON(&actionLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.First(&actionLogModel, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Model(&actionLogModel).Updates(actionLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, actionLogModel)
}

func (h *ActionLogHandler) DeleteActionLog(c *gin.Context) {
	var actionLog models.ActionLog
	if err := h.db.First(&actionLog, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Delete(&actionLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, actionLog)
}
