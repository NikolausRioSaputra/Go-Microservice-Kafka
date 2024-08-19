package handler

import (
	"log"
	"net/http"
	"service-orchestration/m/internal/repository"
	"service-orchestration/m/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	repo        repository.OcresRepository
	kafkaWriter repository.KafkaWriterRepository
}

func NewTransactionHandler(repo repository.OcresRepository, kafkaWriter repository.KafkaWriterRepository) *TransactionHandler {
	return &TransactionHandler{
		repo:        repo,
		kafkaWriter: kafkaWriter,
	}
}

func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.repo.GetAllTransactions()
	if err != nil {
		log.Printf("Error getting transactions: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateItemIdAndResend(c *gin.Context) {
	transactionId := c.Param("transactionId")
	newItemId := c.Query("newItemId")

	if newItemId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New ItemId is required"})
		return
	}

	transactionUseCase := usecase.NewTransactionUseCase(h.repo, h.kafkaWriter)
	err := transactionUseCase.UpdateItemIdAndSendKafkaMessage(transactionId, newItemId)
	if err != nil {
		log.Printf("Error updating and resending transaction: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update and resend transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated and message resent"})
}

func (h *TransactionHandler) UpdatePaymentAndResend(c *gin.Context) {
    transactionId := c.Param("transactionId")
    newPayment := c.Query("paymentMethod")

    if newPayment == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "New PaymentMethod is required"})
        return
    }

    transactionUseCase := usecase.NewTransactionUseCase(h.repo, h.kafkaWriter)
    err := transactionUseCase.UpdatePaymentAndSendKafkaMessage(transactionId, newPayment)
    if err != nil {
        log.Printf("Error updating and resending transaction: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update and resend transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Transaction updated and message resent"})
}
