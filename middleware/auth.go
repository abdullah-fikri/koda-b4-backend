package middleware

import (
	"backend/lib"
	"backend/models"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
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
	frontEnd := os.Getenv("FRONTEND")
    return cors.New(cors.Config{
        AllowOrigins: []string{
            frontEnd,
            "http://localhost:5173",
        },
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })
}
