package handlers

import (
	"net/http"
	"strconv"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupHandler struct {
	db *gorm.DB
}

func NewGroupHandler(db *config.Config) *GroupHandler {
	return &GroupHandler{db: db.DB}
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
	var groups []models.Groups
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * limit
	query := h.db
	if err := query.Limit(limit).Offset(offset).Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.Groups{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}
	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, groups, meta)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	var group models.Group
	id := c.Param("id")
	if err := h.db.Preload("UserGroups.User").Preload("UserGroups.User.UserRoles.Role").First(&group, id).Error; err != nil {
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
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, group)
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()
	var group models.Groups

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.CreateGroup,
		Details: "Create group: " + group.Name,
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, group)
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	var group models.UpdateGroupRequest
	var groupModel models.Groups

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID

	id := c.Param("id")
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.First(&groupModel, id).Error; err != nil {
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
	if err := h.db.Model(&groupModel).Updates(group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	group.ID = groupModel.ID
	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.CreateGroup,
		Details: "Create group: " + group.Name,
	}
	services.CreateActionLog(actionLog, h.db)
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, group)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()
	id := c.Param("id")
	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	if err := tx.Delete(&models.Groups{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteGroup,
		Details: "Delete group: " + id,
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, nil)
}
