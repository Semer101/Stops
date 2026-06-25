package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/models"
	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

const (
	ContextUserID = "userID"
	ContextRole   = "role"
)

func RequireAuth(tokens *services.TokenService) gin.HandlerFunc {
	return func(context *gin.Context) {
		rawHeader := context.GetHeader("Authorization")
		if rawHeader == "" || !strings.HasPrefix(rawHeader, "Bearer ") {
			abortWithError(context, utils.ErrUnauthorized)
			return
		}

		claims, err := tokens.Parse(strings.TrimPrefix(rawHeader, "Bearer "))
		if err != nil {
			abortWithError(context, err)
			return
		}

		context.Set(ContextUserID, claims.UserID)
		context.Set(ContextRole, claims.Role)
		context.Next()
	}
}

func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(context *gin.Context) {
		currentRole, exists := context.Get(ContextRole)
		if !exists || currentRole != role {
			abortWithError(context, utils.ErrForbidden)
			return
		}

		context.Next()
	}
}

func abortWithError(context *gin.Context, err error) {
	context.AbortWithStatusJSON(utils.ErrorStatus(err), utils.ErrorResponse{
		Success: false,
		Message: utils.PublicMessage(err),
		Error:   err.Error(),
	})
}
