package middleware

import (
	"backend/lib"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(401, models.Response{
				Success: false,
				Message: "Unauthorized",
			})
			ctx.Abort()
			return
		}

		payload := user.(lib.UserPayload)

		if payload.Role != "admin" {
			ctx.JSON(403, models.Response{
				Success: false,
				Message: "Forbidden: Admin only",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
