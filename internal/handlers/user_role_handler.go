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

type UserRoleHandler struct {
	db *gorm.DB
}

func NewUserRoleHandler(cfg *config.Config) *UserRoleHandler {
	return &UserRoleHandler{db: cfg.DB}
}

func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {
	var userRoles []models.UserRole
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
	query := h.db.Preload("User").Preload("Role")
	if err := query.Limit(limit).Offset(offset).Find(&userRoles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.UserRole{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}
	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, userRoles, meta)
}

func (h *UserRoleHandler) GetUserRole(c *gin.Context) {
	var userRole models.UserRole

	id := c.Param("id")

	if id == "" {
		reponseFail := map[string]string{
			"name": utils.StatusParameterIsRequired,
		}
		utils.RespondJSON(c, http.StatusCreated, utils.StatusParameterIsRequired, reponseFail)
		return
	}

	if err := h.db.Preload("User").Preload("Role").First(&userRole, id).Error; err != nil {
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
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, userRole)
}

func (h *UserRoleHandler) CreateUserRoles(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()
	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var userRole models.CreateUserRoleRequest

	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole.AssignedAt = time.Now()

	if err := tx.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Create user role %d with user_id %d", userRole.RoleID, userRole.UserID),
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, userRole)
}

func (h *UserRoleHandler) UpdateUserRoles(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID

	var updateUserRoleRequest models.UpdateUserRoleRequest
	var userRole models.UserRole

	id := c.Param("id")

	if id == "" {
		reponseFail := map[string]string{
			"name": utils.StatusParameterIsRequired,
		}
		utils.RespondJSON(c, http.StatusCreated, utils.StatusParameterIsRequired, reponseFail)
		return
	}

	if err := h.db.First(&userRole, id).Error; err != nil {
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

	if err := c.ShouldBindJSON(&updateUserRoleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Model(&userRole).Updates(updateUserRoleRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Update user role %d with user_id %d", userRole.RoleID, userRole.UserID),
	}

	services.CreateActionLog(actionLog, tx)
	tx.Commit()

	utils.RespondJSON(c, http.StatusAccepted, utils.StatusSuccessfully, updateUserRoleRequest)
}

func (h *UserRoleHandler) DeleteUserRoles(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var userRole models.UserRole

	id := c.Param("id")

	if id == "" {
		reponseFail := map[string]string{
			"name": utils.StatusParameterIsRequired,
		}
		utils.RespondJSON(c, http.StatusCreated, utils.StatusParameterIsRequired, reponseFail)
		return
	}

	if err := h.db.First(&userRole, id).Error; err != nil {
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

	if err := tx.Delete(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Delete user role %d with user_id %d", userRole.RoleID, userRole.UserID),
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, userRole)
}
