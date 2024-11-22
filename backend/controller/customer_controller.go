package controller

import (
	"net/http"

	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/service"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	customerService *service.CustomerService
}

func NewCustomerController(customerService *service.CustomerService) *CustomerController {
	return &CustomerController{
		customerService: customerService,
	}
}

func (c *CustomerController) GetAllCustomers(ctx *gin.Context) {
	customers, err := c.customerService.GetAllCustomers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customers)
}

// func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
// 	customerID := ctx.Param("id")

// 	var customer models.Customer
// 	if err := ctx.ShouldBindJSON(&customer); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := c.customerService.UpdateCustomer(ctx, customerID, &customer); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
// }

func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	customerID := ctx.Param("id")
	if customerID == "" {
		ctx.JSON(http.StatusOK, gin.H{"error": "customer id not found"})
		return
	}

	var customer models.Customer
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error(), "message": "error binding customer data", "success": false})
		return
	}

	// Convert to update map and check if there are fields to update
	updates := customer.ToUpdateMapCustomer()
	if len(updates) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "no fields to update", "success": false})
		return
	}

	if err := c.customerService.UpdateCustomer(ctx, customerID, updates); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error(), "message": "failed to update customer", "success": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "success": true})
}
