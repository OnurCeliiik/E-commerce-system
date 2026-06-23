package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/OnurCeliiik/ecommerce/services/inventory/dto"
	"github.com/OnurCeliiik/ecommerce/services/inventory/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryService interface {
	GetInventory(ctx context.Context, productID uuid.UUID) (*dto.InventoryResponse, error)
	UpdateInventory(ctx context.Context, productID uuid.UUID, req dto.UpdateInventoryRequest) (*dto.InventoryResponse, error)
}

type InventoryHandler struct {
	inventoryService InventoryService
}

func NewInventoryHandler(inventoryService InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	resp, err := h.inventoryService.GetInventory(c.Request.Context(), productID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInventoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": service.ErrInventoryNotFound.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.inventoryService.UpdateInventory(c.Request.Context(), productID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrInvalidInput.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
