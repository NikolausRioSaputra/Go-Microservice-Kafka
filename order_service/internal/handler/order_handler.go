package handler

import (
	"net/http"
	"order-service/internal/domain"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	useCase usecase.OrderUseCase
}

func NewOrderHandler(useCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{useCase: useCase}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderRequest domain.OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderRequest.TransactionID = uuid.New().String()

	err := h.useCase.ProcessOrder(c.Request.Context(), orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process order"})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"status": "Order processed successfully"})
	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
		"order":   orderRequest,
	})
}

func (h *OrderHandler) RegisterEvent(c *gin.Context) {
    var registrationRequest domain.EventRegistrationRequest
    if err := c.ShouldBindJSON(&registrationRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    registrationRequest.TransactionID = uuid.New().String()

    err := h.useCase.ProcessEventRegistration(c.Request.Context(), registrationRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register event"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Event registration successful",
        "registration": registrationRequest,
    })
}
