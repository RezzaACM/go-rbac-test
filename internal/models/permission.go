package models

type Permission struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	BaseModel
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func (u *CreatePermissionRequest) TableName() string {
	return "permissions"
}

func (u *UpdatePermissionRequest) TableName() string {
	return "permissions"
}
