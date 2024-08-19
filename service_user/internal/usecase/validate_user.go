package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"user_service/internal/domain"

	retryit "github.com/benebobaa/retry-it"
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

	var apiResponse struct {
		IsValid bool   `json:"isValid"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		fmt.Println("Attempt: ", counter)
		return requestUserValidation(ctx, apiUrl, msg.UserId, &apiResponse)
	}, retryit.WithInitialDelay(1*time.Second), retryit.WithMaxAttempts(5))

	if err != nil {
		return domain.Response{}, fmt.Errorf("error making request: %w", err)
	}

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
			Amount:        msg.Amount,
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
		Amount:        msg.Amount,
	}, nil
}

func requestUserValidation(ctx context.Context, url, userId string, response any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}

	q := req.URL.Query()
	q.Add("userId", userId)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status -> ", resp.StatusCode)
	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	// Process successful response here
	fmt.Printf("Request successful with status code: %d\n", resp.StatusCode)
	return nil
}
