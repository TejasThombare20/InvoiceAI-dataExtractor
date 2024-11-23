package routes

import (
	"github.com/TejasThombare20/backend/controller"
	"github.com/TejasThombare20/backend/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(routes *gin.Engine) {

	extractionService := service.NewExtractionService()
	extractionController := controller.NewExtractionController(extractionService)

	productService := service.NewProductService()
	productController := controller.NewProductController(productService)

	customerService := service.NewCustomerService()
	customerController := controller.NewCustomerController(customerService)

	invoiceService := service.NewInvoiceService()
	invoiceController := controller.NewInvoiceController(invoiceService)

	// API routes
	routes.POST("/extract", extractionController.ExtractData)

	routes.GET("/products", productController.GetAllProducts)
	routes.PUT("/product/:id", productController.UpdateProduct)

	routes.GET("/customers", customerController.GetAllCustomers)
	routes.PUT("/customer/:id", customerController.UpdateCustomer)

	routes.GET("/invoices", invoiceController.GetAllInvoices)
	routes.PUT("/invoice/:id", invoiceController.UpdateInvoice)

}
