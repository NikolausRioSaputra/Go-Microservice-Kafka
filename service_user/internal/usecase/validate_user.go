package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"user_service/internal/domain"
)

type MessageUseCase interface {
	ValidateUser(ctx context.Context, msg domain.Message) (domain.Response, error)
}

type messageUseCase struct{}

func NewMessageUseCase() MessageUseCase {
	return &messageUseCase{}
}

func (uc *messageUseCase) ValidateUser(ctx context.Context, msg domain.Message) (domain.Response, error) {

	apiUrl := "https://uservalidation.free.beeceptor.com"
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return domain.Response{}, err
	}
	// Tambahkan query parameters atau headers jika diperlukan
	q := req.URL.Query()
	q.Add("userId", msg.UserId)
	req.URL.RawQuery = q.Encode()

	// Atur timeout dan buat HTTP client
	client := &http.Client{Timeout: 5 * time.Second}

	// Panggil API eksternal
	resp, err := client.Do(req)
	if err != nil {
		return domain.Response{}, err
	}
	defer resp.Body.Close()
	// Proses response dari API eksternal
	if resp.StatusCode != http.StatusOK {
		return domain.Response{}, errors.New("failed to validate user")
	}

	var apiResponse struct {
		IsValid bool   `json:"isValid"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.Response{}, err
	}

	// Business logic to validate the user
	// Sesuaikan response berdasarkan hasil validasi
	if !apiResponse.IsValid {
		return domain.Response{
			OrderType:     msg.OrderType,
			OrderService:  "validateUser",
			OderID:        msg.OderID,
			TransactionId: msg.TransactionId,
			UserId:        msg.UserId,
			ItemId:        msg.ItemId,
			OrderAmount:   msg.OrderAmount,
			PaymentMethod: msg.PaymentMethod,
			RespCode:      400,
			RespStatus:    apiResponse.Status,
			RespMessage:   apiResponse.Message,
		}, nil
	}

	return domain.Response{
		OrderType:     msg.OrderType,
		OrderService:  "validateUser",
		OderID:        msg.OderID,
		TransactionId: msg.TransactionId,
		UserId:        msg.UserId,
		ItemId:        msg.ItemId,
		OrderAmount:   msg.OrderAmount,
		PaymentMethod: msg.PaymentMethod,
		RespCode:      200,
		RespStatus:    apiResponse.Status,
		RespMessage:   apiResponse.Message,
	}, nil
}
