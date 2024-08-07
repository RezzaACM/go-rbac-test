package services

import (
	"example.com/go-crud/internal/models"
	"gorm.io/gorm"
)

func ValidateRoleName(h *gorm.DB, role models.UpdateRoleRequest, id string) bool {
	var existingRole models.Role
	if err := h.Where("name = ? AND id != ?", role.Name, id).First(&existingRole).Error; err == nil {
		return true
	}
	return false
}
