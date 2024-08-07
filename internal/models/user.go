package models

type User struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	BaseModel
}

type CreateUserRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLoginResponse struct {
	User
	Token string `json:"token"`
}

func (u *CreateUserRequest) TableName() string {
	return "users"
}
