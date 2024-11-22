package service

import (
	"context"

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
