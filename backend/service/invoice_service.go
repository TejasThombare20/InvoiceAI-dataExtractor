package service

import (
	"context"
	"fmt"

	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/repository"
)

type InvoiceService struct {
	repo *repository.InvoiceRepository
}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{
		repo: repository.NewInvoiceRepository(),
	}
}

func (s *InvoiceService) GetAllInvoices(ctx context.Context) ([]models.Invoice, error) {
	return s.repo.FindAll(ctx)
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, id string, updates map[string]interface{}) error {
	// Add validation logic if needed
	// For example, validate serial number format, check total amount is positive, etc.
	if updates["totalAmount"] != nil {
		totalAmount := updates["totalAmount"].(float64)
		if totalAmount < 0 {
			return fmt.Errorf("total amount cannot be negative")
		}
	}

	return s.repo.Update(ctx, id, updates)
}
