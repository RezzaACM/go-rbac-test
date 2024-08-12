package models

import "time"

type UserGroup struct {
	ID         int       `json:"id"`
	GroupID    int       `json:"-"`
	Group      Groups    `json:"group" gorm:"foreignKey:group_id"`
	UserID     int       `json:"-"`
	User       User      `json:"user" gorm:"foreignKey:user_id"`
	AssignedAt time.Time `json:"assigned_at"`
	BaseModel
}

type GroupUser struct {
	ID         int       `json:"id"`
	GroupID    int       `json:"-"`
	UserID     int       `json:"-"`
	User       User      `json:"user" gorm:"foreignKey:user_id"`
	AssignedAt time.Time `json:"assigned_at"`
	BaseModel
}

type CreateUserGroupRequest struct {
	GroupID    int       `json:"group_id"`
	UserID     int       `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

type UpdateUserGroupRequest struct {
	ID         int       `json:"id"`
	GroupID    int       `json:"group_id"`
	UserID     int       `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

func (u *UserGroup) TableName() string {
	return "user_groups"
}

func (u *GroupUser) TableName() string {
	return "user_groups"
}

// TableName returns the database table name for the CreateUserGroup struct.
//
// No parameters.
// Returns a string representing the database table name.
func (u *CreateUserGroupRequest) TableName() string {
	return "user_groups"
}

// TableName returns the database table name for the UpdateUserGroupRequest struct.
//
// No parameters.
// Returns a string representing the database table name.
func (u *UpdateUserGroupRequest) TableName() string {
	return "user_groups"
}
