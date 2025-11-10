package controllers

import (
	"backend/lib"
	"backend/models"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
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

	if err := ctx.BindJSON(&req); err != nil {
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
	var req models.RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	u, err := models.UpdateUser(req.Email, req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "update succesfully",
		Data:    u,
	})
}

func ForgotUSer(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, err := models.Forgot(req.Email)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "email not found",
		})
		return
	}

	intId := int(user.ID)
	token := lib.GeneratedTokens(intId, user.Role)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "account reset success",
		Data: map[string]any{
			"user":  user,
			"token": token,
		},
	})
}

func UploadPicture(ctx *gin.Context) {
	id := ctx.Param("id")

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

	filename := "profile-picture-" + id + ext
	path := "./uploads/" + filename

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal menyimpan file",
		})
		return
	}

	if err := models.UpdateProfilePicture(id, path); err != nil {
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

//admin

func RegisterAd(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, err := models.RegisterAd(req)
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

func UpdateUserAd(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	u, err := models.UpdateUser(req.Email, req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "update succesfully",
		Data:    u,
	})
}
