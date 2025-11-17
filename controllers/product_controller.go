package controllers

import (
	"backend/config"
	"backend/lib"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// admin
func AdminProductList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	search := ctx.DefaultQuery("search", "")

	products, total, err := models.GetProductsAdmin(page, limit, search)
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "admin product list",
		Data: gin.H{
			"products": products,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}

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
    search := ctx.DefaultQuery("q", "")
    sort := ctx.DefaultQuery("sort", "")
    minPrice := ctx.DefaultQuery("min_price", "")
    maxPrice := ctx.DefaultQuery("max_price", "")

    categoryStr := ctx.QueryArray("category[]")

    var categoryIDs []int
    for _, c := range categoryStr {
        id, err := strconv.Atoi(c)
        if err == nil {
            categoryIDs = append(categoryIDs, id)
        }
    }

    page, _ := strconv.Atoi(pageStr)
    if page < 1 { page = 1 }
    limit, _ := strconv.Atoi(limitStr)
    if limit < 1 { limit = 10 }

    isCachable := (search == "" && sort == "" && minPrice == "" && maxPrice == "" && len(categoryIDs) == 0)
    var cacheKey string

    if isCachable {
        cacheKey = fmt.Sprintf("products:page:%d:limit:%d", page, limit)
        var cacheData struct {
            Products   []models.Product `json:"products"`
            Pagination map[string]any   `json:"pagination"`
        }

        cache, err := config.Rdb.Get(context.Background(), cacheKey).Result()
        if err == nil && cache != "" {
            _ = json.Unmarshal([]byte(cache), &cacheData)
            ctx.JSON(200, models.Response{
                Success: true,
                Message: "data from cache",
				Pagination: cacheData.Pagination,
                Data:    cacheData.Products,
            })
            return
        }
    }

    products, total, err := models.GetProducts(page, limit, search, sort, parseInt(minPrice), parseInt(maxPrice), categoryIDs)
    if err != nil {
        ctx.JSON(400, models.Response{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    totalPage := int((total + int64(limit) - 1) / int64(limit))

	queryParams := ctx.Request.URL.Query()
	queryParams.Set("limit", strconv.Itoa(limit))
	
	baseURL := "http://hifiy-backend.vercel.app" + ctx.Request.URL.Path
	var nextURL, prevURL string
	if page > 1 {
		qp := ctx.Request.URL.Query()
		qp.Set("limit", strconv.Itoa(limit))
		qp.Set("page", strconv.Itoa(page-1))
		prevURL = baseURL + "?" + qp.Encode()
	}

	if page < totalPage {
		qp := ctx.Request.URL.Query()
		qp.Set("limit", strconv.Itoa(limit))
		qp.Set("page", strconv.Itoa(page+1))
		nextURL = baseURL + "?" + qp.Encode()
	}
	
    pagination := map[string]any{
        "page":       page,
        "limit":      limit,
        "total_data": total,
        "total_page": totalPage,
        "next":       nextURL,
        "prev":       prevURL,
    }

	cacheData := struct {
		Products   []models.Product `json:"products"`
		Pagination map[string]any   `json:"pagination"`
	}{
		Products:   products,
		Pagination: pagination,
	}

    if isCachable {
        jsonData, err := json.Marshal(cacheData)
        if err == nil {
            config.Rdb.Set(context.Background(), cacheKey, jsonData, 15*time.Minute)
        }
    }

    ctx.JSON(200, models.Response{
        Success: true,
        Message: "success data from db",
		Pagination: cacheData.Pagination,
        Data:    cacheData.Products,
    })
}

func parseInt(s string) int {
    v, _ := strconv.Atoi(s)
    return v
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
		c.JSON(400, models.Response{
			Success: false,
			Message: "product not found",
		})
		return
	}

	recommendations, _ := models.GetRecommendationsByCategory(product.Category, product.ID)

	c.JSON(200, models.Response{
		Success: true,
		Message: "success",
		Data: gin.H{
			"product":         product,
			"recommendations": recommendations,
		},
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
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "invalid product id"})
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "file not provided: " + err.Error()})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := []string{".jpg", ".jpeg", ".png"}

	valid := false
	for _, e := range allowedExt {
		if ext == e {
			valid = true
			break
		}
	}
	if !valid {
		ctx.JSON(400, models.Response{Success: false, Message: "invalid file extension. Only .jpg, .jpeg, .png allowed"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: "cannot open file: " + err.Error()})
		return
	}
	defer src.Close()

	// cloudinary
	uploadedURL, err := lib.UploadImage(src)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: "failed upload to cloudinary: " + err.Error()})
		return
	}
	
	if err := models.UploadImgProduct(int64(productID), uploadedURL); err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "upload success",
		Data:    gin.H{"image_url": uploadedURL},
	})
}
