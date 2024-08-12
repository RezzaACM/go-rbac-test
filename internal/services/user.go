package services

import (
	"fmt"
	"time"

	"example.com/go-crud/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = []byte("secret-key")

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares a hashed password with a plain-text password.
func CheckPasswordHash(password, hash string) bool {
	er := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return er == nil
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func GetUserByUsername(username string, db *gorm.DB) (models.UserLoggedIn, error) {
	var user models.User
	var userLoggedIn models.UserLoggedIn

	if err := db.Where("username = ?", username).Preload("UserRoles.Role").First(&user).Error; err != nil {
		return userLoggedIn, err
	}
	userLoggedIn.ID = user.ID
	userLoggedIn.Fullname = user.Fullname
	userLoggedIn.Username = user.Username
	userLoggedIn.Email = user.Email
	userLoggedIn.UserRoles = user.UserRoles
	return userLoggedIn, nil
}
