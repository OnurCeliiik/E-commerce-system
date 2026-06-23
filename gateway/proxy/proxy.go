package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// UserService forwards HTTP requests to the user microservice.
type UserService struct {
	proxy *httputil.ReverseProxy
}

type ProductService struct {
	proxy *httputil.ReverseProxy
}

type InventoryService struct {
	proxy *httputil.ReverseProxy
}

// NewUserService creates a reverse proxy to the given base URL,
// e.g. http://user-service:8080 inside Docker.
func NewUserService(targetURL string) (*UserService, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("parse user service url: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "user service unavailable", http.StatusBadGateway)
	}

	return &UserService{proxy: proxy}, nil
}

func (p *UserService) ServeHTTP(c *gin.Context) {
	p.proxy.ServeHTTP(c.Writer, c.Request)
}

// NewProductService creates a reverse proxy to the given base URL,
// e.g. http://product-service:8080 inside Docker.
func NewProductService(targetURL string) (*ProductService, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("parse product service url: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "product service unavailable", http.StatusBadGateway)
	}
	return &ProductService{proxy: proxy}, nil
}

func (p *ProductService) ServeHTTP(c *gin.Context) {
	p.proxy.ServeHTTP(c.Writer, c.Request)
}

func NewInventoryService(targetURL string) (*InventoryService, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("parse inventory service url: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "inventory service unavailable", http.StatusBadGateway)
	}
	return &InventoryService{proxy: proxy}, nil
}

func (p *InventoryService) ServeHTTP(c *gin.Context) {
	p.proxy.ServeHTTP(c.Writer, c.Request)
}
