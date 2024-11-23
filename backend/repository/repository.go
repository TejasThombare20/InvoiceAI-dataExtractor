package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TejasThombare20/backend/config"
	"github.com/TejasThombare20/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExtractionRepository struct {
	db *mongo.Database
}

func NewExtractionRepository() *ExtractionRepository {
	return &ExtractionRepository{
		db: config.DB,
	}
}

func (r *ExtractionRepository) SaveExtractedData(ctx context.Context, extractedData *models.ExtractedData, filename string) error {
	// Start a new transaction
	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	// Define collections
	invoiceCollection := r.db.Collection("invoices")
	productCollection := r.db.Collection("products")
	customerCollection := r.db.Collection("customers")

	// Perform the transaction
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// Insert Invoice
		invoice := models.Invoice{
			ID:           primitive.NewObjectID(),
			SerialNumber: &extractedData.Invoice.SerialNumber,
			TotalAmount:  &extractedData.Invoice.TotalAmount,
			CreatedAt:    time.Now(),
		}

		fmt.Println("date invoice", extractedData.Invoice.Date)
		parsedDate, err := time.Parse("2 Jan 2006", extractedData.Invoice.Date)
		if err != nil {
			log.Fatalf("Invalid date format: %v", err)
		}

		invoice.Date = &parsedDate
		invoiceResult, err := invoiceCollection.InsertOne(sessionContext, invoice)
		if err != nil {
			return fmt.Errorf("failed to insert invoice: %v", err)
		}
		invoiceID := invoiceResult.InsertedID.(primitive.ObjectID)

		// Insert Products
		var productDocuments []interface{}
		for _, product := range extractedData.Products {
			productDoc := models.Product{
				ID:           primitive.NewObjectID(),
				Name:         product.Name,
				Quantity:     &product.Quantity,
				UnitPrice:    &product.UnitPrice,
				Tax:          &product.Tax,
				PriceWithTax: &product.PriceWithTax,
				InvoiceID:    invoiceID,
			}
			productDocuments = append(productDocuments, productDoc)
		}
		if len(productDocuments) > 0 {
			_, err = productCollection.InsertMany(sessionContext, productDocuments)
			if err != nil {
				return fmt.Errorf("failed to insert products: %v", err)
			}
		}

		// Insert Customer
		customer := models.Customer{
			ID:                  primitive.NewObjectID(),
			Name:                &extractedData.Customer.Name,
			PhoneNumber:         &extractedData.Customer.PhoneNumber,
			TotalPurchaseAmount: &extractedData.Customer.TotalPurchaseAmount,
			InvoiceID:           invoiceID,
		}
		_, err = customerCollection.InsertOne(sessionContext, customer)
		if err != nil {
			return fmt.Errorf("failed to insert customer: %v", err)
		}

		return nil
	})

	return err
}

// Helper function to validate extracted data
func ValidateExtractedData(data *models.ExtractedData) error {
	// Check Invoice details
	// if data.Invoice == nil {
	// 	return fmt.Errorf("missing invoice details")
	// }

	if data.Invoice.SerialNumber == "" {
		return fmt.Errorf("missing invoice serial number")
	}

	if data.Invoice.Date == "" {
		return fmt.Errorf("missing invoice date")
	}

	if data.Invoice.TotalAmount <= 0 {
		return fmt.Errorf("invalid total amount")
	}

	// Check Products
	if len(data.Products) == 0 {
		return fmt.Errorf("no products found")
	}

	for _, product := range data.Products {
		if product.Name == "" {
			return fmt.Errorf("product name is missing")
		}
		if product.Quantity <= 0 {
			return fmt.Errorf("invalid product quantity")
		}
		if product.UnitPrice < 0 {
			return fmt.Errorf("invalid product unit price")
		}
	}

	// Check Customer
	// if data.Customer == nil {
	// 	return fmt.Errorf("missing customer details")
	// }

	if data.Customer.Name == "" {
		return fmt.Errorf("missing customer name")
	}

	if data.Customer.PhoneNumber == "" {
		return fmt.Errorf("missing customer phone number")
	}

	return nil
}
