package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OnurCeliiik/ecommerce/services/product/dto"
	"github.com/OnurCeliiik/ecommerce/services/product/handlers"
	"github.com/OnurCeliiik/ecommerce/services/product/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockProductService struct {
	create       func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	listProducts func(ctx context.Context) ([]dto.ProductResponse, error)
	getByID      func(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
	update       func(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	delete       func(ctx context.Context, id uuid.UUID) error
}

func (m *mockProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	return m.create(ctx, req)
}

func (m *mockProductService) ListProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	return m.listProducts(ctx)
}

func (m *mockProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	return m.getByID(ctx, id)
}

func (m *mockProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	return m.update(ctx, id, req)
}

func (m *mockProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return m.delete(ctx, id)
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCreateProduct_Success(t *testing.T) {
	handler := handlers.NewProductHandler(&mockProductService{
		create: func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
			return &dto.ProductResponse{Name: req.Name, Price: req.Price}, nil
		},
	})

	body := `{"name":"Widget","price":9.99,"description":"test"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateProduct(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	handler := handlers.NewProductHandler(&mockProductService{
		getByID: func(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
			return nil, service.ErrProductNotFound
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/products/"+uuid.New().String(), nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler.GetProductByID(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetProductByID_InvalidID(t *testing.T) {
	handler := handlers.NewProductHandler(&mockProductService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/products/not-a-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	handler.GetProductByID(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	id := uuid.New()
	handler := handlers.NewProductHandler(&mockProductService{
		update: func(ctx context.Context, productID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
			return &dto.ProductResponse{ID: productID, Name: req.Name, Price: req.Price}, nil
		},
	})

	body := `{"name":"Updated Widget","price":12.99,"description":"updated"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/products/"+id.String(), bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.UpdateProduct(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	handler := handlers.NewProductHandler(&mockProductService{
		update: func(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
			return nil, service.ErrProductNotFound
		},
	})

	body := `{"name":"Updated Widget","price":12.99}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/products/"+uuid.New().String(), bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler.UpdateProduct(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateProduct_InvalidID(t *testing.T) {
	handler := handlers.NewProductHandler(&mockProductService{})

	body := `{"name":"Updated Widget","price":12.99}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/products/not-a-uuid", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	handler.UpdateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	id := uuid.New()
	handler := handlers.NewProductHandler(&mockProductService{
		delete: func(ctx context.Context, productID uuid.UUID) error {
			return nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.DeleteProduct(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", w.Body.String())
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	handler := handlers.NewProductHandler(&mockProductService{
		delete: func(ctx context.Context, id uuid.UUID) error {
			return service.ErrProductNotFound
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+uuid.New().String(), nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler.DeleteProduct(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
