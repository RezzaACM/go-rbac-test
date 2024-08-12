package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/go-crud/internal/config"
	userService "example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from the Authorization header.
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondJSON(ctx, http.StatusUnauthorized, utils.StatusLoginHeaderRequired, nil)
			ctx.Abort()
			return
		}

		// Check if the token is in the correct format.
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		if tokenString == "" {
			utils.RespondJSON(ctx, http.StatusUnauthorized, utils.StatusLoginHeaderRequired, nil)
			ctx.Abort()
			return
		}

		// Parse tokenString
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			utils.RespondJSON(ctx, http.StatusUnauthorized, utils.StatusLoginInvalidToken, err)
			ctx.Abort()
			return
		}

		// Validate the token
		valid, err := userService.VerifyToken(tokenString)
		if err != nil || !valid {
			utils.RespondJSON(ctx, http.StatusUnauthorized, utils.StatusLoginInvalidToken, err)
			ctx.Abort()
			return
		}

		// Add the user to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username := claims["username"].(string)
			// var user models.User
			// call connection database
			value, err := userService.GetUserByUsername(username, cfg.DB)

			if err != nil {
				utils.RespondJSON(ctx, http.StatusInternalServerError, utils.StatusSomethingWrong, err)
				ctx.Abort()
				return
			}
			ctx.Set("user", value)
		}

		// Continue to the next handler
		ctx.Next()
	}
}
