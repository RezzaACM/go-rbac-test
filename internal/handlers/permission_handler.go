package handlers

import (
	"net/http"
	"strconv"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type PermissionHandler struct {
	db *gorm.DB
}

func NewPermissionHandler(cfg *config.Config) *PermissionHandler {
	return &PermissionHandler{db: cfg.DB}
}

func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	var permissions []models.Permission
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

	if err := query.Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var count int64
	h.db.Model(&models.Permission{}).Count(&count)
	meta := gin.H{
		"page":       page,
		"limit":      limit,
		"totalPages": (count + int64(limit) - 1) / int64(limit),
		"count":      count,
	}

	utils.ResponseJSONWithMeta(c, http.StatusOK, utils.StatusSuccessfully, permissions, meta)
}

func (h *PermissionHandler) GetPermission(c *gin.Context) {
	var permission models.Permission
	id := c.Param("id")

	if err := h.db.First(&permission, id).Error; err != nil {
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
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, permission)
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var permission models.CreatePermissionRequest

	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	permission.Slug = slug.Make(permission.Name)
	if err := h.db.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, permission)

}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	var permission models.UpdatePermissionRequest
	var getPermisison models.Permission

	id := c.Param("id")

	if id == "" {
		reponseFail := map[string]string{
			"name": utils.StatusParameterIsRequired,
		}
		utils.RespondJSON(c, http.StatusCreated, utils.StatusParameterIsRequired, reponseFail)
		return
	}

	if err := h.db.First(&getPermisison, id).Error; err != nil {
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

	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission.ID = getPermisison.ID
	permission.Slug = slug.Make(permission.Name)
	if err := h.db.Save(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.RespondJSON(c, http.StatusAccepted, utils.StatusSuccessfully, permission)
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	var permission models.Permission
	id := c.Param("id")

	if err := h.db.First(&permission, id).Error; err != nil {
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

	if err := h.db.Delete(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	placeholders := map[string]string{
		"name": permission.Name,
	}

	utils.RespondJSON(c, http.StatusAccepted, utils.ReplacePlaceholders(utils.StatusDeletedDataSucessfully, placeholders), permission)
}
