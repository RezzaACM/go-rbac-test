package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleHandler struct {
	db *gorm.DB
}

func NewRoleHandler(cfg *config.Config) *RoleHandler {
	return &RoleHandler{db: cfg.DB}
}

func (h *RoleHandler) GetRoles(c *gin.Context) {
	var roles []models.Role
	name := c.Query("name")
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

	query := h.db

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var count int64
	h.db.Model(&models.Role{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}

	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, roles, meta)

}

func (h *RoleHandler) GetRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := h.db.First(&role, id).Error; err != nil {
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
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var role models.CreateRoleRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role.Name = strings.ToLower(role.Name)

	var existingRole models.Role
	if err := h.db.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		placeholders := map[string]string{
			"name": role.Name,
		}
		utils.RespondJSON(c, http.StatusForbidden, utils.ReplacePlaceholders(utils.StatusRoleAlreadyUsed, placeholders), role)
		return
	}

	if err := tx.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.CreateRole,
		Details: "Create role: " + role.Name,
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var role models.UpdateRoleRequest
	var fetchRole models.Role
	id := c.Param("id")

	if id == "" {
		reponseFail := map[string]string{
			"name": utils.StatusParameterIsRequired,
		}
		utils.RespondJSON(c, http.StatusCreated, utils.StatusParameterIsRequired, reponseFail)
		return
	}

	if err := h.db.First(&fetchRole, id).Error; err != nil {
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

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if checkRoleName := services.ValidateRoleName(h.db, role, id); checkRoleName {
		placeholders := map[string]string{
			"name": role.Name,
		}
		utils.RespondJSON(c, http.StatusForbidden, utils.ReplacePlaceholders(utils.StatusRoleAlreadyUsed, placeholders), role)
		return
	}
	role.ID = id

	if err := tx.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.UpdateRole,
		Details: "Update role: " + role.Name,
	}
	services.CreateActionLog(actionLog, tx)
	tx.Commit()
	utils.RespondJSON(c, http.StatusAccepted, utils.StatusSuccessfully, role)

}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	// add transaction database
	tx := h.db.Begin()

	// get c.Get user id from middleware
	user, _ := c.Get("user")
	userID := user.(models.UserLoggedIn).ID
	var role models.Role
	id := c.Param("id")

	if err := h.db.First(&role, id).Error; err != nil {
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

	if err := h.db.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	placeholders := map[string]string{
		"name": role.Name,
	}

	// add action log to record input
	actionLog := models.CreateActionLogRequest{
		UserId:  uint(userID),
		Action:  utils.DeleteRole,
		Details: "Delete role: " + role.Name,
	}
	services.CreateActionLog(actionLog, tx)

	tx.Commit()
	utils.RespondJSON(c, http.StatusAccepted, utils.ReplacePlaceholders(utils.StatusDeletedDataSucessfully, placeholders), role)

}
