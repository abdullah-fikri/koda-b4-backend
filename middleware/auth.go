package middleware

import (
	"backend/lib"
	"backend/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(401, models.Response{
				Success: false,
				Message: "Unauthorized",
			})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := lib.VerifyToken(token)
		if err != nil {
			ctx.JSON(401, models.Response{
				Success: false,
				Message: "Token invalid",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", payload)

		ctx.Set("user_id", int64(payload.Id))
		ctx.Set("role", payload.Role)
		ctx.Next()

	}
}

func AllowPreflight(r *gin.Engine) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
			return
		}
		ctx.Next()
	}
}
