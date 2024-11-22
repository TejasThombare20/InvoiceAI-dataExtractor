package controller

import (
	"net/http"

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
