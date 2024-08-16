package domain

// Struct Message digunakan untuk merepresentasikan data pesan yang diterima.
type Message struct {
	OrderType     string  `json:"orderType"`
	OderID        string  `json:"orderID"`
	TransactionId string  `json:"transactionId"`
	UserId        string  `json:"userId"`
	Price         float64 `json:"price"`
	ItemId        string  `json:"itemId"`
	PaymentMethod string  `json:"paymentMethod"`
	OrderAmount   int     `json:"orderAmount"`
}

// Struct Response digunakan untuk merepresentasikan data balasan setelah pesan diproses.
type Response struct {
	OrderType     string  `json:"orderType"`
	OrderService  string  `json:"orderService"`
	OderID        string  `json:"orderID"`
	PaymentMethod string  `json:"paymentMethod"`
	OrderAmount   int     `json:"orderAmount"`
	TransactionId string  `json:"transactionId"`
	UserId        string  `json:"userId"`
	ItemId        string  `json:"itemId"`
	Price         float64 `json:"price"`
	RespCode      int     `json:"respCode"`
	RespStatus    string  `json:"respStatus"`
	RespMessage   string  `json:"respMessage"`
}
