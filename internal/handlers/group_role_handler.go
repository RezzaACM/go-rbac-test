package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupRoleHandler struct {
	db *gorm.DB
}

func NewGroupRoleHandler(db *config.Config) *GroupRoleHandler {
	return &GroupRoleHandler{db: db.DB}
}

func (h *GroupRoleHandler) GetGroupRoles(c *gin.Context) {
	var groupRoles []models.GroupRole
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
	query := h.db.Preload("Group").Preload("Role")
	if err := query.Limit(limit).Offset(offset).Find(&groupRoles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.GroupRole{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}
	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, groupRoles, meta)
}

func (h *GroupRoleHandler) GetGroupRole(c *gin.Context) {
	var groupRole models.GroupRole

	id := c.Param("id")

	if err := h.db.Preload("Group").Preload("Role").First(&groupRole, id).Error; err != nil {
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
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, groupRole)
}

func (h *GroupRoleHandler) CreateGroupRole(c *gin.Context) {
	var groupRole models.CreateGroupRoleRequest
	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	if err := c.ShouldBindJSON(&groupRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Create(&groupRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteGroup,
		Details: fmt.Sprintf("Create group role: %d - %d", groupRole.GroupID, groupRole.RoleID),
	}
	services.CreateActionLog(actionLog, h.db)
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, groupRole)
}

func (h *GroupRoleHandler) UpdateGroupRole(c *gin.Context) {
	var groupRole models.UpdateGroupRoleRequest
	var groupRoleModel models.GroupRole

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID

	id := c.Param("id")
	if err := h.db.First(&groupRoleModel, id).Error; err != nil {
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
	if err := c.ShouldBindJSON(&groupRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Model(&groupRoleModel).Updates(groupRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	groupRole.ID = groupRoleModel.ID

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteGroup,
		Details: fmt.Sprintf("Update group role: %d - %d", groupRole.GroupID, groupRole.RoleID),
	}
	services.CreateActionLog(actionLog, h.db)
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, groupRole)
}

func (h *GroupRoleHandler) DeleteGroupRole(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()
	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var groupRole models.GroupRole
	id := c.Param("id")
	if err := h.db.First(&groupRole, id).Error; err != nil {
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
	if err := tx.Delete(&groupRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteGroup,
		Details: fmt.Sprintf("Delete group role: %d - %d", groupRole.GroupID, groupRole.RoleID),
	}
	services.CreateActionLog(actionLog, h.db)
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, nil)
}
