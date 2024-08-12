package models

import "time"

type RolePermission struct {
	ID           uint       `json:"id"`
	RoleID       uint       `json:"-"`
	Role         Role       `json:"role" gorm:"foreignKey:role_id"`
	PermissionID uint       `json:"-"`
	Permission   Permission `json:"permission" gorm:"foreignKey:permission_id"`
	AssignedAt   time.Time  `json:"assigned_at"`
	BaseModel
}

type CreateRolePermissionRequest struct {
	RoleID       uint      `json:"role_id"`
	PermissionID uint      `json:"permission_id"`
	AssignedAt   time.Time `json:"assigned_at"`
	BaseModel
}

type UpdateRolePermissionRequest struct {
	ID           uint      `json:"id"`
	RoleID       uint      `json:"role_id"`
	PermissionID uint      `json:"permission_id"`
	AssignedAt   time.Time `json:"assigned_at"`
	BaseModel
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (CreateRolePermissionRequest) TableName() string {
	return "role_permissions"
}

func (UpdateRolePermissionRequest) TableName() string {
	return "role_permissions"
}
