package models

type Groups struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseModel
}

type Group struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	UserGroups  []GroupUser `json:"user_group"`
	BaseModel
}

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateGroupRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *Groups) TableName() string {
	return "groups"
}

func (h *Group) TableName() string {
	return "groups"
}

func (h *CreateGroupRequest) TableName() string {
	return "groups"
}

func (h *UpdateGroupRequest) TableName() string {
	return "groups"
}
