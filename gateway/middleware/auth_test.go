package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OnurCeliiik/ecommerce/gateway/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockTokenValidator struct {
	userID uuid.UUID
	role   string
	err    error
}

func (m *mockTokenValidator) UserIDFromToken(token string) (uuid.UUID, error) {
	if m.err != nil {
		return uuid.Nil, m.err
	}
	return m.userID, nil
}

func (m *mockTokenValidator) RoleFromToken(token string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.role, nil
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuth_MissingHeader(t *testing.T) {
	router := gin.New()
	router.GET("/protected", middleware.Auth(&mockTokenValidator{}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	router := gin.New()
	router.GET("/protected", middleware.Auth(&mockTokenValidator{userID: uuid.New(), role: "customer"}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequireRole_ForbiddenForCustomer(t *testing.T) {
	router := gin.New()
	router.POST("/products",
		middleware.Auth(&mockTokenValidator{userID: uuid.New(), role: "customer"}),
		middleware.RequireRole(middleware.RoleAdmin),
		func(c *gin.Context) {
			c.Status(http.StatusCreated)
		},
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/products", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRequireRole_AllowsAdmin(t *testing.T) {
	router := gin.New()
	router.POST("/products",
		middleware.Auth(&mockTokenValidator{userID: uuid.New(), role: middleware.RoleAdmin}),
		middleware.RequireRole(middleware.RoleAdmin),
		func(c *gin.Context) {
			c.Status(http.StatusCreated)
		},
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/products", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
