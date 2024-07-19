package handlers

import (
	"net/http"
	"strings"

	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/models"
	userService "example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{db: cfg.DB}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondJSON(c, http.StatusOK, utils.StatusSuccessfully, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := h.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		utils.RespondJSON(c, http.StatusBadRequest, utils.StatusEmailAlreadyUsed, user)
		return
	}

	if err := h.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		utils.RespondJSON(c, http.StatusBadRequest, utils.StatusUsernameAlreadyUsed, user)
		return
	}

	hashedPassword, err := userService.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = hashedPassword
	user.Email = strings.ToLower(user.Email)

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondJSON(c, http.StatusCreated, utils.StatusSuccessfully, nil)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Email = strings.ToLower(request.Email)

	var user models.User
	if err := h.db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		utils.RespondJSON(c, http.StatusUnauthorized, "Email is invalid", nil)
		return
	}

	if !userService.CheckPasswordHash(request.Password, user.Password) {
		utils.RespondJSON(c, http.StatusUnauthorized, "Password is invalid", nil)
		return
	}

	jwtToken, err := userService.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseReturn := models.UserLoginResponse{
		User:  user,
		Token: jwtToken,
	}

	utils.RespondJSON(c, http.StatusOK, utils.StatusLoginIsSuccessfully, responseReturn)
}
