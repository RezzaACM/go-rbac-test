package models

import (
	"time"

	"gorm.io/gorm"
)

type GroupRole struct {
	ID         uint      `json:"id"`
	GroupID    uint      `json:"-"`
	Group      Groups    `json:"group" gorm:"foreignKey:group_id"`
	RoleID     uint      `json:"-"`
	Role       Role      `json:"role" gorm:"foreignKey:role_id"`
	AssignedAt time.Time `json:"assigned_at"`
	BaseModel
}

type CreateGroupRoleRequest struct {
	GroupID uint `json:"group_id"`
	RoleID  uint `json:"role_id"`
}

type UpdateGroupRoleRequest struct {
	ID      uint `json:"id"`
	GroupID uint `json:"group_id"`
	RoleID  uint `json:"role_id"`
}

func (h *GroupRole) TableName() string {
	return "group_roles"
}
func (h *GroupRole) BeforeCreate(tx *gorm.DB) (err error) {
	h.AssignedAt = time.Now()

	return
}

func (h *CreateGroupRoleRequest) TableName() string {
	return "group_roles"
}

func (h *UpdateGroupRoleRequest) TableName() string {
	return "group_roles"
}
