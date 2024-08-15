package domain

type Message struct {
	OrderType     string  `json:"orderType"`
	TransactionId string  `json:"transactionId"`
	UserId        string  `json:"userId"`
	ItemId        string  `json:"itemId"`
	OrderAmount   float64 `json:"orderAmount"`
	PaymentMethod string  `json:"paymentMethod"`
}

type Response struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	ItemId        string `json:"itemId"`
	RespCode      int    `json:"respCode"`
	RespStatus    string `json:"respStatus"`
	RespMessage   string `json:"respMessage"`
}
