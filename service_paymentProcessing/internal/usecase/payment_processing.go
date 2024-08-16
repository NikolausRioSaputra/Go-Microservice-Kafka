package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"service_paymentProcessing/internal/domain"
	"time"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, msg domain.PaymentMessage) (domain.PaymentResponse, error)
}

type paymentUseCase struct {
}

func NewPaymentUseCase() PaymentUseCase {
	return &paymentUseCase{}
}

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, msg domain.PaymentMessage) (domain.PaymentResponse, error) {

	apiUrl := "https://paymentprocessing.free.beeceptor.com"
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return domain.PaymentResponse{}, err
	}
	// Tambahkan query parameters atau headers jika diperlukan
	q := req.URL.Query()
	q.Add("paymentMethod", msg.PaymentMethod)
	req.URL.RawQuery = q.Encode()

	// Atur timeout dan buat HTTP client
	client := &http.Client{Timeout: 5 * time.Second}

	// Panggil API eksternal
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return domain.PaymentResponse{}, err
	}
	defer resp.Body.Close()

	// Proses response dari API eksternal
	if resp.StatusCode != http.StatusOK {
		return domain.PaymentResponse{}, errors.New("failed to process payment")
	}

	var apiResponse struct {
		// IsSuccess bool    `json:"isSuccess"`
		Balance   float64 `json:"balance"`
		Status    string  `json:"status"`
		Message   string  `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Println(err.Error())
		return domain.PaymentResponse{}, err
	}

	totalPrice := float64(msg.OrderAmount) * msg.Price

	if apiResponse.Balance < totalPrice {
		return domain.PaymentResponse{
			OrderType:     msg.OrderType,
			OrderService:  "processPayment",
			TransactionId: msg.TransactionId,
			OrderID:       msg.OrderID,
			UserId:        msg.UserId,
			Balance:       apiResponse.Balance,
			PaymentMethod: msg.PaymentMethod,
			OrderAmount:   msg.OrderAmount,
			Price:         msg.Price,
			ItemId:        msg.ItemId,
			RespCode:      400,
			RespStatus:    apiResponse.Status,
			RespMessage:   apiResponse.Message,
		}, nil
	}

	return domain.PaymentResponse{
		OrderType:     msg.OrderType,
		OrderService:  "processPayment",
		OrderID:       msg.OrderID,
		TransactionId: msg.TransactionId,
		UserId:        msg.UserId,
		Balance:       apiResponse.Balance - totalPrice,
		Price:         msg.Price,
		OrderAmount:   msg.OrderAmount,
		PaymentMethod: msg.PaymentMethod,
		ItemId:        msg.ItemId,
		RespCode:      200,
		RespStatus:    apiResponse.Status,
		RespMessage:   apiResponse.Message,
	}, nil
}
