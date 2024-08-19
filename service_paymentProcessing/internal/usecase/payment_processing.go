package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"service_paymentProcessing/internal/domain"
	"time"

	retryit "github.com/benebobaa/retry-it"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, msg domain.Message) (domain.Response, error)
}

type messageUseCase struct {
}

func NewPaymentUseCase() PaymentUseCase {
	return &messageUseCase{}
}

func (uc *messageUseCase) ProcessPayment(ctx context.Context, msg domain.Message) (domain.Response, error) {
	apiUrl := "https://paymentprocessing.free.beeceptor.com"

	var apiResponse struct {
		Balance float64 `json:"balance"`
		Status  string  `json:"status"`
		Message string  `json:"message"`
	}

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		fmt.Println("Payment Processing Attempt: ", counter)
		return requestProcessPayment(ctx, apiUrl, msg.PaymentMethod, &apiResponse)
	}, retryit.WithInitialDelay(2*time.Second), retryit.WithMaxAttempts(5))

	if err != nil {
		return domain.Response{}, fmt.Errorf("error processing payment: %w", err)
	}

	var totalPrice float64
	if msg.OrderType == "Register Event" {
		totalPrice = float64(msg.Amount)
	}

	if msg.OrderType == "Buy Item" {
		// Calculate total price based on order amount and price
		totalPrice = float64(msg.OrderAmount) * msg.Price
	}

	log.Println("total price: ", totalPrice)
	log.Println("amount: ", msg.Amount)
	if apiResponse.Balance < totalPrice {
		return domain.Response{
			OrderType:     msg.OrderType,
			OrderService:  "processPayment",
			OrderID:       msg.OrderID,
			TransactionId: msg.TransactionId,
			UserId:        msg.UserId,
			Balance:       apiResponse.Balance,
			PaymentMethod: msg.PaymentMethod,
			OrderAmount:   msg.OrderAmount,
			Amount:        msg.Amount,
			Price:         msg.Price,
			ItemId:        msg.ItemId,
			RespCode:      400,
			RespStatus:    "Failed",
			RespMessage:   "Insufficient balance",
		}, nil
	}

	return domain.Response{
		OrderType:     msg.OrderType,
		OrderService:  "processPayment",
		OrderID:       msg.OrderID,
		TransactionId: msg.TransactionId,
		UserId:        msg.UserId,
		Balance:       apiResponse.Balance - totalPrice,
		Amount:        msg.Amount,
		Price:         msg.Price,
		OrderAmount:   msg.OrderAmount,
		PaymentMethod: msg.PaymentMethod,
		ItemId:        msg.ItemId,
		RespCode:      200,
		RespStatus:    "Success",
		RespMessage:   apiResponse.Message,
	}, nil
}

func requestProcessPayment(ctx context.Context, url, paymentMethod string, response any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}

	q := req.URL.Query()
	q.Add("paymentMethod", paymentMethod)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("Payment Processing Response Status -> ", resp.StatusCode)
	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	fmt.Printf("Payment processing request successful with status code: %d\n", resp.StatusCode)
	return nil
}
