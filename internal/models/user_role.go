package models

import "time"

type UserRole struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"-"`
	RoleID     uint      `json:"-"`
	User       User      `json:"user" gorm:"foreignKey:user_id"`
	Role       Role      `json:"role" gorm:"foreignKey:role_id"`
	AssignedAt time.Time `json:"assigned_at"`
	BaseModel
}

type UserRolesUser struct {
	UserID uint `json:"-"`
	RoleID uint `json:"-"`
	Role   Role `json:"role" gorm:"foreignKey:role_id"`
}

type CreateUserRoleRequest struct {
	UserID     uint      `json:"user_id"`
	RoleID     uint      `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

type UpdateUserRoleRequest struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	RoleID     uint      `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

func (u *CreateUserRoleRequest) TableName() string {
	return "user_roles"
}
func (u *UpdateUserRoleRequest) TableName() string {
	return "user_roles"
}

func (u *UserRolesUser) TableName() string {
	return "user_roles"
}
