package domain

// Struct Message digunakan untuk merepresentasikan data pesan yang diterima.
type Message struct {
	OrderType     string  `json:"orderType"`
	OrderService  string  `json:"orderService"`
	TransactionId string  `json:"transactionId"`
	OrderID       string  `json:"orderID"`
	UserId        string  `json:"userId"`
	Price         float64 `json:"price"`
	ItemId        string  `json:"itemId"`
	PaymentMethod string  `json:"paymentMethod"`
	OrderAmount   int     `json:"orderAmount"`
	Amount        int     `json:"amount"`
}

// Struct Message digunakan untuk merepresentasikan data pesan yang diterima.
type Response struct {
	OrderType     string  `json:"orderType"`
	OrderService  string  `json:"orderService"`
	Price         float64 `json:"price"`
	TransactionId string  `json:"transactionId"`
	OrderID       string  `json:"orderID"`
	PaymentMethod string  `json:"paymentMethod"`
	Balance       float64 `json:"balance"`
	UserId        string  `json:"userId"`
	OrderAmount   int     `json:"orderAmount"`
	ItemId        string  `json:"itemId"`
	RespCode      int     `json:"respCode"`
	RespStatus    string  `json:"respStatus"`
	RespMessage   string  `json:"respMessage"`
	Amount        int     `json:"amount"`
}
