package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Product godoc
// @Summary Get all products
// @Description Public product list
// @Tags Products
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} models.Response
// @Router /products [get]
func Product(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")
	sort := ctx.DefaultQuery("sort", "")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	key := ctx.Request.RequestURI
	var cacheData struct {
		Products   []models.Product `json:"products"`
		Pagination map[string]any   `json:"pagination"`
	}

	cache, err := config.Rdb.Get(context.Background(), key).Result()
	if err == nil && cache != "" {
		_ = json.Unmarshal([]byte(cache), &cacheData)

		ctx.JSON(200, models.Response{
			Success: true,
			Message: "data from cache",
			Data:    cacheData,
		})
		return
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

	cacheData = struct {
		Products   []models.Product `json:"products"`
		Pagination map[string]any   `json:"pagination"`
	}{
		Products: products,
		Pagination: map[string]any{
			"page":       page,
			"limit":      limit,
			"total_data": total,
			"total_page": totalPage,
		},
	}

	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "failed marshal data",
		})
		return
	}
	config.Rdb.Set(context.Background(), key, jsonData, 15*time.Minute)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "success data from db",
		Data:    cacheData,
	})
}

// ProductDetail godoc
// @Summary Get product detail by ID
// @Description Fetch product detail
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Router /products/{id} [get]
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
		Data:    product,
	})
}

// CreateProduct godoc
// @Summary Create new product
// @Description Admin creates a new product
// @Tags Admin - Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.CreateProductRequest true "Product data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /admin/product [post]
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

	redisCtx := context.Background()
	iter := config.Rdb.Scan(redisCtx, 0, "/products*", 0).Iterator()
	for iter.Next(redisCtx) {
		config.Rdb.Del(redisCtx, iter.Val())
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Product created",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary Update product
// @Description Admin updates product data
// @Tags Admin - Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param product body models.UpdateProductRequest true "Updated product data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /admin/product/{id} [put]
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

	redisCtx := context.Background()
	iter := config.Rdb.Scan(redisCtx, 0, "/products*", 0).Iterator()
	for iter.Next(redisCtx) {
		config.Rdb.Del(redisCtx, iter.Val())
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product updated",
		Data:    product,
	})
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Admin can delete product
// @Tags Products (Admin)
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products/{id} [delete]
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

	redisCtx := context.Background()
	iter := config.Rdb.Scan(redisCtx, 0, "/products*", 0).Iterator()
	for iter.Next(redisCtx) {
		config.Rdb.Del(redisCtx, iter.Val())
	}
	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product deleted",
	})
}

// UploadProductImages godoc
// @Summary Upload product image
// @Description Upload image for product (.jpg, .jpeg, .png max 10MB)
// @Tags Products (Admin)
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param image formData file true "Product Image"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products/{id}/image [post]
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
