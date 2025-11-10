package controllers

import (
	"backend/lib"
	"backend/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	roleVal, exists := ctx.Get("role")

	if !exists {
		req.Role = "user"
	} else {
		role := roleVal.(string)
		if role == "admin" {
			req.Role = "user"
		} else {
			req.Role = "user"
		}
	}

	user, err := models.Register(req)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Register success",
		Data:    user,
	})
}

func LoginUser(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, err := models.Login(req.Email)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "wrong email or password",
		})
		return
	}

	if !lib.VerifyPassword(req.Password, user.Password) {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "wrong email or password",
		})
		return
	}
	intId := int(user.ID)
	token := lib.GeneratedTokens(intId, user.Role)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Login success",
		Data: map[string]any{
			"user":  user,
			"token": token,
		},
	})

}
func UpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	targetID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "invalid user id"})
		return
	}

	userID := ctx.MustGet("user_id").(int64)
	role := ctx.MustGet("role").(string)

	if role != "admin" && userID != targetID {
		ctx.JSON(403, models.Response{
			Success: false,
			Message: "you cannot update another user's data",
		})
		return
	}

	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	userEmail, err := models.GetUserEmailByID(targetID)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "user not found",
		})
		return
	}

	updated, err := models.UpdateUser(userEmail, req)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "update succesfully",
		Data:    updated,
	})
}

func UploadPicture(ctx *gin.Context) {
	idParam := ctx.Param("id")

	paramUserID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "invalid user id",
		})
		return
	}

	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var userID int64
	switch v := userIDInterface.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "invalid user id format",
		})
		return
	}

	userRole := ctx.MustGet("role").(string)

	if userRole != "admin" && userID != paramUserID {
		ctx.JSON(403, models.Response{
			Success: false,
			Message: "you cannot change another user's profile picture",
		})
		return
	}

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "file not found",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := []string{".jpg", ".jpeg", ".png"}
	valid := false
	for _, v := range allowed {
		if ext == v {
			valid = true
			break
		}
	}
	if !valid {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "format harus .jpg .jpeg .png",
		})
		return
	}

	if file.Size > 10<<20 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "file maksimal 10MB",
		})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal membaca file",
		})
		return
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	_, err = openedFile.Read(buffer)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal validasi file",
		})
		return
	}

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "file bukan gambar valid",
		})
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal membuat direktori upload",
		})
		return
	}

	filename := fmt.Sprintf("profile-picture-%d-%d%s", paramUserID, time.Now().Unix(), ext)
	path := filepath.Join(uploadDir, filename)

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal menyimpan file",
		})
		return
	}

	if err := models.UpdateProfilePicture(paramUserID, path); err != nil {
		os.Remove(path)
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Upload success",
		Data: map[string]any{
			"profile_picture": path,
		},
	})
}


func ListUser(ctx *gin.Context){
	user, err := models.ListUser()
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "list all user",
		Data: user,
	})
}