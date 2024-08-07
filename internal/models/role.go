package models

type Role struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    int    `json:"parent_id"`
	BaseModel
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
}

type UpdateRoleRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
}

func (u *CreateRoleRequest) TableName() string {
	return "roles"
}

func (u *UpdateRoleRequest) TableName() string {
	return "roles"
}
