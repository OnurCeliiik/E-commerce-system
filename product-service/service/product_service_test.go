package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/product/dto"
	"github.com/OnurCeliiik/ecommerce/services/product/model"
	"github.com/OnurCeliiik/ecommerce/services/product/repository"
	"github.com/OnurCeliiik/ecommerce/services/product/service"
	"github.com/google/uuid"
)

type mockProductRepository struct {
	create   func(ctx context.Context, product *model.Product) error
	list     func(ctx context.Context) ([]model.Product, error)
	findByID func(ctx context.Context, id uuid.UUID) (*model.Product, error)
}

func (m *mockProductRepository) Create(ctx context.Context, product *model.Product) error {
	return m.create(ctx, product)
}

func (m *mockProductRepository) List(ctx context.Context) ([]model.Product, error) {
	return m.list(ctx)
}

func (m *mockProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	return m.findByID(ctx, id)
}

func TestCreateProduct_Success(t *testing.T) {
	repo := &mockProductRepository{
		create: func(ctx context.Context, product *model.Product) error {
			product.CreatedAt = time.Now()
			return nil
		},
	}

	svc := service.NewProductService(repo)
	resp, err := svc.CreateProduct(context.Background(), dto.CreateProductRequest{
		Name:        "Widget",
		Price:       9.99,
		Description: "A widget",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Name != "Widget" {
		t.Fatalf("expected Widget, got %s", resp.Name)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uuid.UUID) (*model.Product, error) {
			return nil, repository.ErrProductNotFound
		},
	}

	svc := service.NewProductService(repo)
	_, err := svc.GetProductByID(context.Background(), uuid.New())
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestListProducts_Success(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepository{
		list: func(ctx context.Context) ([]model.Product, error) {
			return []model.Product{{
				ID:        id,
				Name:      "Widget",
				Price:     9.99,
				CreatedAt: time.Now(),
			}}, nil
		},
	}

	svc := service.NewProductService(repo)
	products, err := svc.ListProducts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 || products[0].Name != "Widget" {
		t.Fatalf("unexpected products: %+v", products)
	}
}
