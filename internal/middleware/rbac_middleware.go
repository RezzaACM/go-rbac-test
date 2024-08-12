package middleware

import (
	"net/http"

	"example.com/go-crud/internal/models"
	"example.com/go-crud/internal/utils"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exist := ctx.Get("user")
		if !exist {
			utils.RespondJSON(ctx, http.StatusUnauthorized, utils.StatusLoginHeaderRequired, nil)
			ctx.Abort()
			return
		}

		userLoggedIn := user.(models.UserLoggedIn)
		if len(userLoggedIn.UserRoles) == 0 {
			utils.RespondJSON(ctx, http.StatusForbidden, utils.StatusMismatchRole, nil)
			ctx.Abort()
			return
		}

		// check if the user has any of the allowed roles
		for _, role := range allowedRoles {
			if userLoggedIn.UserRoles[0].Role.Name == role {
				ctx.Next()
				return
			}
		}
		utils.RespondJSON(ctx, http.StatusForbidden, utils.StatusMismatchRole, nil)
		ctx.Abort()
	}
}
