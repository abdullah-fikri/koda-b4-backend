package controllers

import (
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

func Product(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")
	sort := ctx.DefaultQuery("sort", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	products, total, err := models.GetProducts(page, limit, search, sort)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "products fetched",
		Data: map[string]any{
			"products": products,
			"pagination": map[string]any{
				"page":       page,
				"limit":      limit,
				"total_data": total,
				"total_page": totalPage,
			},
		},
	})
}

func ProductDetail(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product id",
		})
		return
	}

	product, err := models.GetProductByID(productID)
	if err != nil {
		c.JSON(400,
			models.Response{
				Success: false,
				Message: "product not found",
			})
		return
	}

	c.JSON(200, models.Response{
		Success: true,
		Message: "success",
		Data: product,
	})
}

func CreateProduct(ctx *gin.Context) {
	var req models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	product, err := models.CreateProduct(req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Product created",
		Data:    product,
	})
}

func UpdateProduct(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	product, err := models.UpdateProduct(int64(id), req)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product updated",
		Data:    product,
	})
}
func DeleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product id",
		})
		return
	}

	err = models.DeleteProduct(int64(id))
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product deleted",
	})
}

func UploadProductImages(ctx *gin.Context) {
	productIDParam := ctx.Param("id")
	productID, err := strconv.ParseInt(productIDParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "invalid product id",
		})
		return
	}

	role := ctx.MustGet("role").(string)
	if role != "admin" {
		ctx.JSON(403, models.Response{
			Success: false,
			Message: "only admin can upload product images",
		})
		return
	}

	file, err := ctx.FormFile("image")
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

	uploadDir := "./product_uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal membuat direktori upload",
		})
		return
	}

	filename := fmt.Sprintf("product-%d-%d%s", productID, time.Now().Unix(), ext)
	path := filepath.Join(uploadDir, filename)

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "gagal menyimpan file",
		})
		return
	}

	err = models.UploadImgProduct(productID, path)
	if err != nil {
		os.Remove(path)
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Upload product image success",
		Data: map[string]any{
			"image": path,
		},
	})
}
