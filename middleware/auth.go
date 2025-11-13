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
func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "http://coffe-shop-one-eta.vercel.app")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH,PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200) 
			return
		}
		
		ctx.Next()
	}
}
