package controllers

import (
	"backend/config"
	"backend/lib"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /auth/register [post]
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

	config.Rdb.Del(context.Background(), "/users")
	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Register success",
		Data:    user,
	})
}

// LoginUser godoc
// @Summary Login user
// @Description Login with email and password to get token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /auth/login [post]
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

// UpdateUser godoc
// @Summary Update user data
// @Description User can update their account, admin can update any account
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.RegisterRequest true "Updated Data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Router /user/{id} [put]
func AdminUpdateUser(ctx *gin.Context) {
    idParam := ctx.Param("id")
    targetID, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        ctx.JSON(400, models.Response{Success: false, Message: "invalid user id"})
        return
    }

    role := ctx.MustGet("role").(string)
    if role != "admin" {
        ctx.JSON(403, models.Response{
            Success: false,
            Message: "only admin can update user",
        })
        return
    }

    var req models.AdminUpdateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
        return
    }

    updated, err := models.AdminUpdateUserByID(targetID, req)
    if err != nil {
        ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
        return
    }

    config.Rdb.Del(ctx, "/users")

    ctx.JSON(200, models.Response{
        Success: true,
        Message: "admin update user successfully",
        Data:    updated,
    })
}

//user
func UploadUserPicture(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int64)

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "file not provided"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := []string{".jpg", ".jpeg", ".png"}
	if !slices.Contains(allowed, ext) {
		ctx.JSON(400, models.Response{Success: false, Message: "invalid file format"})
		return
	}

	newFilename := fmt.Sprintf("%d-%d%s", userID, time.Now().Unix(), ext)
	uploadPath := "./uploads/profile/" + newFilename
	os.MkdirAll("./uploads/profile", 0755)

	if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: "failed to save file"})
		return
	}

	if err := models.UpdateUserProfilePicture(userID, newFilename); err != nil {
		os.Remove(uploadPath)
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "upload success",
		Data: gin.H{"profile_picture": newFilename},
	})
}

//admin
func AdminUploadUserPicture(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)
	if role != "admin" {
		ctx.JSON(403, models.Response{Success: false, Message: "forbidden"})
		return
	}

	param := ctx.Param("id")
	targetUserID, err := strconv.ParseInt(param, 10, 64)
	if err != nil || targetUserID <= 0 {
		ctx.JSON(400, models.Response{Success: false, Message: "invalid user id"})
		return
	}

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "file not provided"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := []string{".jpg", ".jpeg", ".png"}
	if !slices.Contains(allowed, ext) {
		ctx.JSON(400, models.Response{Success: false, Message: "invalid file format"})
		return
	}

	newFilename := fmt.Sprintf("%d-%d%s", targetUserID, time.Now().Unix(), ext)
	uploadPath := "./uploads/profile/" + newFilename
	os.MkdirAll("./uploads/profile", 0755)

	if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: "failed to save file"})
		return
	}

	if err := models.AdminUpdateUserProfilePicture(targetUserID, newFilename); err != nil {
		os.Remove(uploadPath)
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "admin updated user picture successfully",
		Data: gin.H{"profile_picture": newFilename},
	})
}


// ListUser godoc
// @Summary Get all users
// @Description Only admin can view all user list
// @Tags Admin
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /admin/users [get]
func ListUser(ctx *gin.Context) {
	key := ctx.Request.RequestURI
	var users []models.ListUserStruct

	cache, err := config.Rdb.Get(context.Background(), key).Result()
	if err != nil || cache == "" {
		users, err = models.ListUser()
		if err != nil {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		data, _ := json.Marshal(users)
		config.Rdb.Set(context.Background(), key, data, 15*time.Second)

		ctx.JSON(200, models.Response{
			Success: true,
			Message: "success data from db",
			Data:    users,
		})
		return
	}

	json.Unmarshal([]byte(cache), &users)
	ctx.JSON(200, models.Response{
		Success: true,
		Message: "data from cache",
		Data:    users,
	})
}


func UserProfile(ctx *gin.Context){
	userData,_ := ctx.Get("user")
	user := userData.(lib.UserPayload)

	profile, err := models.GetUserProfile(int64(user.Id))
	if err != nil{
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, models.Response{
		Success: true,
		Message: "success get user profile",
		Data: profile,
	})
}

func UpdateProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int64)

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updated, err := models.UpdateUserByID(userID, req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "profile updated successfully",
		Data:    updated,
	})
}
