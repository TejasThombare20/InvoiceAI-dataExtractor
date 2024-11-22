package service

import (
	"context"
	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		repo: repository.NewProductRepository(),
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.FindAll(ctx)
}

// func (s *ProductService) UpdateProduct(ctx context.Context, id string, product *models.Product) error {

// 	return s.repo.Update(ctx, id, product)
// }

func (s *ProductService) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.repo.Update(ctx, id, updates)
}
