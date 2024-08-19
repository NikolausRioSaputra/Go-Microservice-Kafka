package domain

// struct ini di gunakan untuk menangani pesan kafka yang masuk
type Message struct {
	OrderType     string  `json:"orderType"`
	EventName     string  `json:"eventName"`
	OrderService  string  `json:"orderService,omitempty"`
	OrderID       string  `json:"orderID"`
	Balance       float64 `json:"balance"`
	TransactionId string  `json:"transactionId"`
	PaymentMethod string  `json:"paymentMethod"`
	OrderAmount   int     `json:"orderAmount"`
	Price         float64 `json:"price"`
	UserId        string  `json:"userId"`
	ItemId        string  `json:"itemId"`
	RespCode      int     `json:"respCode,omitempty"`
	RespStatus    string  `json:"respStatus,omitempty"`
	RespMessage   string  `json:"respMessage,omitempty"`
	Payload       string  `json:"payload"`
	Amount        int     `json:"amount"`
}
