package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"service_paymentProcessing/internal/domain"
	"time"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, msg domain.PaymentMessage) (domain.PaymentResponse, error)
}

type paymentUseCase struct{}

func NewPaymentUseCase() PaymentUseCase {
	return &paymentUseCase{}
}

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, msg domain.PaymentMessage) (domain.PaymentResponse, error) {

	apiUrl := "https://paymentprocessing.free.beeceptor.com"
	_, err := json.Marshal(msg)
	if err != nil {
		return domain.PaymentResponse{}, err
	}

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return domain.PaymentResponse{}, err
	}

	// Atur timeout dan buat HTTP client
	client := &http.Client{Timeout: 5 * time.Second}

	// Panggil API eksternal
	resp, err := client.Do(req)
	if err != nil {
		return domain.PaymentResponse{}, err
	}
	defer resp.Body.Close()

	// Proses response dari API eksternal
	if resp.StatusCode != http.StatusOK {
		return domain.PaymentResponse{}, errors.New("failed to process payment")
	}

	var apiResponse struct {
		IsSuccess bool   `json:"isSuccess"`
		Status    string `json:"status"`
		Message   string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.PaymentResponse{}, err
	}

	if !apiResponse.IsSuccess {
		return domain.PaymentResponse{
			OrderType:     msg.OrderType,
			OrderService:  "processPayment",
			TransactionId: msg.TransactionId,
			UserId:        msg.UserId,
			ItemId:        msg.ItemId,
			RespCode:      400,
			RespStatus:    "Failed",
			RespMessage:   "Payment failed",
		}, nil
	}

	return domain.PaymentResponse{
		OrderType:     msg.OrderType,
		OrderService:  "processPayment",
		TransactionId: msg.TransactionId,
		UserId:        msg.UserId,
		ItemId:        msg.ItemId,
		RespCode:      200,
		RespStatus:    "Success",
		RespMessage:   "Payment processed successfully",
	}, nil
}
