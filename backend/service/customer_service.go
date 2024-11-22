package service

import (
	"context"

	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService() *CustomerService {
	return &CustomerService{
		repo: repository.NewCustomerRepository(),
	}
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	return s.repo.FindAll(ctx)
}

// func (s *CustomerService) UpdateCustomer(ctx context.Context, id string, customer *models.Customer) error {
// 	// Add validation logic if needed
// 	return s.repo.Update(ctx, id, customer)
// }

func (s *CustomerService) UpdateCustomer(ctx context.Context, id string, updates map[string]interface{}) error {
	// Add any validation logic here if needed
	return s.repo.Update(ctx, id, updates)
}
