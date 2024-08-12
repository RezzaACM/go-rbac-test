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

type UserGroupHandler struct {
	db *gorm.DB
}

func NewUserGroupHandler(db *config.Config) *UserGroupHandler {
	return &UserGroupHandler{db: db.DB}
}

func (h *UserGroupHandler) GetUserGroups(c *gin.Context) {
	var userGroups []models.UserGroup
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
	query := h.db.Preload("User").Preload("Group")
	if err := query.Limit(limit).Offset(offset).Find(&userGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.UserGroup{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}
	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, userGroups, meta)
}

func (h *UserGroupHandler) GetUserGroup(c *gin.Context) {
	var userGroup models.UserGroup
	id := c.Param("id")

	if err := h.db.Preload("User").Preload("Group").First(&userGroup, id).Error; err != nil {
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
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, userGroup)

}

func (h *UserGroupHandler) CreateUserGroup(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var userGroup models.CreateUserGroupRequest
	if err := c.ShouldBindJSON(&userGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userGroup.AssignedAt = time.Now()
	if err := tx.Create(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Create user group_id %d with user_id %d", userGroup.GroupID, userGroup.UserID),
	}
	services.CreateActionLog(actionLog, tx)

	// commit
	tx.Commit()
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, userGroup)
}

func (h *UserGroupHandler) UpdateUserGroup(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var userGroup models.UpdateUserGroupRequest
	var userGroupModel models.UserGroup
	id := c.Param("id")
	if err := c.ShouldBindJSON(&userGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.First(&userGroupModel, id).Error; err != nil {
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
	if err := tx.Model(&userGroupModel).Updates(userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userGroup.ID = userGroupModel.ID

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Update user group_id %d with user_id %d", userGroup.GroupID, userGroup.UserID),
	}
	services.CreateActionLog(actionLog, tx)

	// commit
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, userGroupModel)
}

func (h *UserGroupHandler) DeleteUserGroup(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	id := c.Param("id")
	if err := h.db.Delete(&models.UserGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: fmt.Sprintf("Delete user group_id %s with user_id %d", id, userID),
	}
	services.CreateActionLog(actionLog, tx)

	// commit
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, nil)
}
