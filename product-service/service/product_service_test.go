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
	update   func(ctx context.Context, product *model.Product) error
	delete   func(ctx context.Context, id uuid.UUID) error
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

func (m *mockProductRepository) Update(ctx context.Context, product *model.Product) error {
	return m.update(ctx, product)
}

func (m *mockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.delete(ctx, id)
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

func TestUpdateProduct_Success(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, productID uuid.UUID) (*model.Product, error) {
			if productID != id {
				t.Fatalf("unexpected id %s", productID)
			}
			return &model.Product{
				ID:          id,
				Name:        "Widget",
				Price:       9.99,
				Description: "A widget",
			}, nil
		},
		update: func(ctx context.Context, product *model.Product) error {
			if product.Name != "Updated Widget" || product.Price != 12.99 {
				t.Fatalf("unexpected product: %+v", product)
			}
			return nil
		},
	}

	svc := service.NewProductService(repo)
	resp, err := svc.UpdateProduct(context.Background(), id, dto.UpdateProductRequest{
		Name:  "Updated Widget",
		Price: 12.99,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Name != "Updated Widget" || resp.Price != 12.99 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	repo := &mockProductRepository{
		findByID: func(ctx context.Context, id uuid.UUID) (*model.Product, error) {
			return nil, repository.ErrProductNotFound
		},
	}

	svc := service.NewProductService(repo)
	_, err := svc.UpdateProduct(context.Background(), uuid.New(), dto.UpdateProductRequest{
		Name: "Updated Widget",
	})
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	id := uuid.New()
	repo := &mockProductRepository{
		delete: func(ctx context.Context, productID uuid.UUID) error {
			if productID != id {
				t.Fatalf("unexpected id %s", productID)
			}
			return nil
		},
	}

	svc := service.NewProductService(repo)
	if err := svc.DeleteProduct(context.Background(), id); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	repo := &mockProductRepository{
		delete: func(ctx context.Context, id uuid.UUID) error {
			return repository.ErrProductNotFound
		},
	}

	svc := service.NewProductService(repo)
	err := svc.DeleteProduct(context.Background(), uuid.New())
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}
