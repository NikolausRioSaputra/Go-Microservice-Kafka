package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	// "errors"
	"net/http"
	"service-package/internal/domain"
	"time"

	retryit "github.com/benebobaa/retry-it"
)

// Interface ini mendefinisikan kontrak untuk MessageUseCase. Di sini, metode ActivatePackage harus diimplementasikan oleh struct yang mengimplementasi interface ini.
type MessageUseCase interface {
	CheckItem(ctx context.Context, msg domain.Message) (domain.Response, error)
}

type messageUseCase struct{}

func NewMessageUseCase() MessageUseCase {
	return &messageUseCase{}
}

func (uc *messageUseCase) CheckItem(ctx context.Context, msg domain.Message) (domain.Response, error) {
	// Business logic to activate the package
	apiUrl := "https://packageactivate.free.beeceptor.com"

	var apiResponse struct {
		IsValid   bool    `json:"isValid"`
		ItemId    string  `json:"itemId"`
		ItemName  string  `json:"itemName"`
		ItemPrice float64 `json:"itemPrice"`
	}

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		fmt.Println("Attempt: ", counter)
		return requestCheckItem(ctx, apiUrl, msg.ItemId, &apiResponse)
	}, retryit.WithInitialDelay(2*time.Second), retryit.WithMaxAttempts(5))

	if err != nil {
		return domain.Response{}, fmt.Errorf("error making request: %w", err)
	}

	if !apiResponse.IsValid {
		return domain.Response{
			OrderType:     msg.OrderType,
			OrderService:  "validateItem",
			OderID:        msg.OderID,
			TransactionId: msg.TransactionId,
			UserId:        msg.UserId,
			ItemId:        msg.ItemId,
			Price:         apiResponse.ItemPrice,
			OrderAmount:   msg.OrderAmount,
			PaymentMethod: msg.PaymentMethod,
			RespCode:      400,
			RespStatus:    "Failed",
			RespMessage:   apiResponse.ItemName + " is not available",
		}, nil
	}

	return domain.Response{
		OrderType:     msg.OrderType,
		OrderService:  "validateItem",
		OderID:        msg.OderID,
		TransactionId: msg.TransactionId,
		UserId:        msg.UserId,
		ItemId:        msg.ItemId,
		Price:         apiResponse.ItemPrice,
		OrderAmount:   msg.OrderAmount,
		PaymentMethod: msg.PaymentMethod,
		RespCode:      200,
		RespStatus:    "Success",
		RespMessage:   apiResponse.ItemName + " is available",
	}, nil

}

func requestCheckItem(ctx context.Context, url, itemId string, response any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}

	q := req.URL.Query()
	q.Add("itemId", itemId)
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
