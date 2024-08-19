package domain

// struct ini di gunakan untuk kirim pesan
type Message struct {
	OrderType     string `json:"orderType"`
	EventName     string `json:"eventName"`
	OrderService  string `json:"orderService,omitempty"`
	OderID        string `json:"orderID"`
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"paymentMethod"`
	OrderAmount   int    `json:"orderAmount"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	ItemId        string `json:"itemId"`
	RespCode      int    `json:"respCode,omitempty"`
	RespStatus    string `json:"respStatus,omitempty"`
	RespMessage   string `json:"respMessage,omitempty"`
}
