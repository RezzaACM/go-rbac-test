package middleware

import (
	"net/http"
	"strings"

	userService "example.com/go-crud/internal/services"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
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

		// Validate the token
		valid, err := userService.VerifyToken(tokenString)
		if err != nil || !valid {
			utils.RespondJSON(ctx, http.StatusUnauthorized, utils.StatusLoginInvalidToken, nil)
			ctx.Abort()
			return
		}

		// Continue to the next handler
		ctx.Next()
	}
}
