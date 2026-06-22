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
	err    error
}

func (m *mockTokenValidator) UserIDFromToken(token string) (uuid.UUID, error) {
	return m.userID, m.err
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
	router.GET("/protected", middleware.Auth(&mockTokenValidator{userID: uuid.New()}), func(c *gin.Context) {
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
