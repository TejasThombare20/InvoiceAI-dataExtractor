package controller

import (
	"net/http"

	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/service"
	"github.com/gin-gonic/gin"
)

type InvoiceController struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceController(invoiceService *service.InvoiceService) *InvoiceController {
	return &InvoiceController{
		invoiceService: invoiceService,
	}
}

func (c *InvoiceController) GetAllInvoices(ctx *gin.Context) {
	// Optional: Add query parameters for filtering, sorting, pagination
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	invoices, err := c.invoiceService.GetAllInvoices(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page":     page,
		"limit":    limit,
		"invoices": invoices,
	})
}

func (c *InvoiceController) UpdateInvoice(ctx *gin.Context) {
	invoiceID := ctx.Param("id")
	if invoiceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invoice id not found",
		})
		return
	}

	var invoice models.Invoice
	if err := ctx.ShouldBindJSON(&invoice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Convert to update map and check if there are fields to update
	updates := invoice.ToUpdateMapInvoice()
	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "no fields to update",
		})
		return
	}

	if err := c.invoiceService.UpdateInvoice(ctx, invoiceID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "failed to update invoice",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Invoice updated successfully",
	})
}
