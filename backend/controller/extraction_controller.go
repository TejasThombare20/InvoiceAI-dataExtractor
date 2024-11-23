package controller

import (
	"fmt"
	"net/http"

	"github.com/TejasThombare20/backend/service"
	"github.com/gin-gonic/gin"
)

type ExtractionController struct {
	extractionService *service.ExtractionService
}

func NewExtractionController(extractionService *service.ExtractionService) *ExtractionController {
	return &ExtractionController{
		extractionService: extractionService,
	}
}

func (c *ExtractionController) ExtractData(ctx *gin.Context) {

	file, header, err := ctx.Request.FormFile("file")

	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "something went wrong!!",
			"error":   "file not found",
		})
		return
	}
	defer file.Close()

	extractedData, err := c.extractionService.ExtractDataFromFile(ctx, file, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "Failed to extract data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "data extracted successfully",
		"data":    extractedData,
	})
}
