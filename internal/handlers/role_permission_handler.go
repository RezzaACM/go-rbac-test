package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RolePermissionHandler struct {
	db *gorm.DB
}

func NewRolePermissionHandler(cfg *config.Config) *RolePermissionHandler {
	return &RolePermissionHandler{db: cfg.DB}
}

func (h *RolePermissionHandler) GetRolePermissions(c *gin.Context) {
	var rolePermissions []models.RolePermission
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
	query := h.db.Preload("Role").Preload("Permission")
	if err := query.Limit(limit).Offset(offset).Find(&rolePermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.RolePermission{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}
	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, rolePermissions, meta)
}

func (h *RolePermissionHandler) GetRolePermission(c *gin.Context) {
	var rolePermission models.RolePermission

	if err := h.db.First(&rolePermission, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			placeholders := map[string]string{
				"name": c.Param("id"),
			}
			utils.RespondJSON(c, http.StatusForbidden, utils.ReplacePlaceholders(utils.StatusDataNotFound, placeholders), nil)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, rolePermission)
}

func (h *RolePermissionHandler) CreateRolePermission(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var rolePermission models.CreateRolePermissionRequest

	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rolePermission.AssignedAt = time.Now()

	if err := tx.Create(&rolePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Created role permission: %d for role: %d", rolePermission.RoleID, rolePermission.PermissionID),
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()

	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, rolePermission)
}

func (h *RolePermissionHandler) UpdateRolePermission(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var rolePermission models.UpdateRolePermissionRequest
	var rolePermissionModel models.RolePermission

	id := c.Param("id")

	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.First(&rolePermissionModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			placeholders := map[string]string{
				"name": id,
			}
			utils.RespondJSON(c, http.StatusForbidden, utils.ReplacePlaceholders(utils.StatusDataNotFound, placeholders), nil)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rolePermission.AssignedAt = time.Now()

	if err := tx.Model(&rolePermissionModel).Updates(rolePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rolePermission.ID = rolePermissionModel.ID

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Updated role permission: %d for role: %d", rolePermission.RoleID, rolePermission.PermissionID),
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, rolePermission)
}

func (h *RolePermissionHandler) DeleteRolePermission(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var rolePermission models.RolePermission

	if err := h.db.First(&rolePermission, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			placeholders := map[string]string{
				"name": c.Param("id"),
			}
			utils.RespondJSON(c, http.StatusForbidden, utils.ReplacePlaceholders(utils.StatusDataNotFound, placeholders), nil)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if err := tx.Delete(&rolePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Deleted role permission: %d for role: %d", rolePermission.RoleID, rolePermission.PermissionID),
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, nil)
}
