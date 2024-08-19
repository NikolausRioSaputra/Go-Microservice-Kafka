package domain

type EventRegistrationRequest struct {
	OrderID       string `json:"orderID"`
	EventName     string `json:"eventName"`
	UserID        string `json:"userId"`
	OrderType     string `json:"orderType"`
	TransactionID string `json:"transactionId"`
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"paymentMethod"`
}

type EventResponse struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService"`
	OrderID       string `json:"orderID"`
	TransactionID string `json:"transactionId"`
	UserID        string `json:"userId"`
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"paymentMethod"`
	RespCode      int    `json:"respCode"`
	RespStatus    string `json:"respStatus"`
	RespMessage   string `json:"respMessage"`
}
