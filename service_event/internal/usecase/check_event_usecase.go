package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-package/internal/domain"
	"time"

	retryit "github.com/benebobaa/retry-it"
)

type EventUseCase interface {
	ProcessEventRegistration(ctx context.Context, req domain.EventRegistrationRequest) (domain.EventResponse, error)
}

type eventUseCase struct{}

func NewEventUseCase() EventUseCase {
	return &eventUseCase{}
}

func (uc *eventUseCase) ProcessEventRegistration(ctx context.Context, req domain.EventRegistrationRequest) (domain.EventResponse, error) {
	// API URL for event validation (replace with your actual API endpoint)
	apiURL := "https://eventvalidation.free.beeceptor.com"

	var apiResponse struct {
		IsValid   bool   `json:"isValid"`
		EventName string `json:"eventName"`
	}

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		fmt.Println("Attempt: ", counter)
		return requestEventValidation(ctx, apiURL, req.EventName, &apiResponse)
	}, retryit.WithInitialDelay(2*time.Second), retryit.WithMaxAttempts(2))

	if err != nil {
		return domain.EventResponse{}, fmt.Errorf("error validating event: %w", err)
	}

	if !apiResponse.IsValid {
		return domain.EventResponse{
			OrderType:     "Register Event",
			OrderService:  "validateEvent",
			OrderID:       req.OrderID,
			TransactionID: req.TransactionID,
			UserID:        req.UserID,
			PaymentMethod: req.PaymentMethod,
			RespCode:      400,
			RespStatus:    "Failed",
			RespMessage:   fmt.Sprintf("Event '%s' is not available", req.EventName),
			Amount:        req.Amount,
		}, nil
	}

	// If the event is valid, return a success response
	return domain.EventResponse{
		OrderType:     "Register Event",
		OrderService:  "validateEvent",
		OrderID:       req.OrderID,
		TransactionID: req.TransactionID,
		UserID:        req.UserID,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		RespCode:      200,
		RespStatus:    "success",
		RespMessage:   "success register event",
	}, nil
}

func requestEventValidation(ctx context.Context, url, eventName string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}

	q := req.URL.Query()
	q.Add("eventName", eventName)
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

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	fmt.Printf("Request successful with status code: %d\n", resp.StatusCode)
	return nil
}
