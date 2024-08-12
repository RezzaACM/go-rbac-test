package models

import (
	"time"

	"gorm.io/gorm"
)

type ActionLog struct {
	Id        uint      `json:"id"`
	UserId    uint      `json:"-"`
	User      User      `json:"user" gorm:"foreignKey:user_id"`
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	BaseModel
}

type CreateActionLogRequest struct {
	UserId  uint   `json:"user_id"`
	Action  string `json:"action"`
	Details string `json:"details"`
}

type UpdateActionLogRequest struct {
	Id      uint   `json:"id"`
	UserId  uint   `json:"user_id"`
	Action  string `json:"action"`
	Details string `json:"details"`
}

func (h *ActionLog) TableName() string {
	return "audit_logs"
}

func (h *CreateActionLogRequest) TableName() string {
	return "audit_logs"
}

func (h *UpdateActionLogRequest) TableName() string {
	return "audit_logs"
}

func (h *ActionLog) BeforeCreate(tx *gorm.DB) (err error) {
	h.Timestamp = time.Now()
	return
}

func (h *ActionLog) BeforeUpdate(tx *gorm.DB) (err error) {
	h.Timestamp = time.Now()
	return
}
