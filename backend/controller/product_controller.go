package controller

import (
	"net/http"

	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/service"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.productService.GetAllProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// func (c *ProductController) UpdateProduct(ctx *gin.Context) {
// 	productID := ctx.Param("id")

// 	if productID == "" {

// 		ctx.JSON(http.StatusOK, gin.H{"message": "product id not found", "success": false})
// 		return
// 	}

// 	var product models.Product
// 	if err := ctx.ShouldBindJSON(&product); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	fmt.Println("update", product)

// 	if err := c.productService.UpdateProduct(ctx, productID, &product); err != nil {

// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "falied to update product", "success": false})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "success": true})
// }

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.JSON(http.StatusOK, gin.H{"message": "product id not found", "success": false})
		return
	}

	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// Convert struct to update map
	updates := product.ToUpdateMapProduct()
	if len(updates) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "no fields to update", "success": false})
		return
	}

	if err := c.productService.UpdateProduct(ctx, productID, updates); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error":   err.Error(),
			"message": "failed to update product",
			"success": false,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "success": true})
}
