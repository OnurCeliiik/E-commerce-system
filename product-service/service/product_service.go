package service

import (
	"context"
	"errors"

	"github.com/OnurCeliiik/ecommerce/services/product/dto"
	"github.com/OnurCeliiik/ecommerce/services/product/model"
	"github.com/OnurCeliiik/ecommerce/services/product/repository"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	List(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *productService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &model.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *productService) ListProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		resp = append(resp, *toProductResponse(&products[i]))
	}
	return resp, nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Stock > 0 {
		product.Stock = req.Stock
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

func toProductResponse(product *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
	}
}
