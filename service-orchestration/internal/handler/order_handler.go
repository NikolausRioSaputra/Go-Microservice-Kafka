package handler

import (
	"context"
	"net/http"
	"service-orchestration/m/internal/domain"
	"service-orchestration/m/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase usecase.OrderUseCase // ->  jembatan antara HTTP request dan logika bisnis (use case).
}

func NewOrderHandler(uc usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: uc,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderReq domain.OrderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.orderUseCase.ProcessOrder(context.Background(), orderReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process order" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": orderReq})
}
