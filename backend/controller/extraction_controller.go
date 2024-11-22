package controller

import (
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "No file provided",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	extractedData, err := c.extractionService.ExtractDataFromFile(ctx, file, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to extract data",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   extractedData,
	})
}
