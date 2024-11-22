package repository

import (
	"context"

	"github.com/TejasThombare20/backend/config"
	"github.com/TejasThombare20/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceRepository struct {
	collection *mongo.Collection
}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{
		collection: config.DB.Collection("invoices"),
	}
}

func (r *InvoiceRepository) FindAll(ctx context.Context) ([]models.Invoice, error) {
	// Add optional filtering and pagination
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invoices []models.Invoice
	if err = cursor.All(ctx, &invoices); err != nil {
		return nil, err
	}

	return invoices, nil
}
