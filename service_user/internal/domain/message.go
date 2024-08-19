package domain

type Message struct {
	OrderType     string `json:"orderType"`
	OderID        string `json:"orderID"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	ItemId        string `json:"itemId"`
	PaymentMethod string `json:"paymentMethod"`
	OrderAmount   int    `json:"orderAmount"`
	Amount        int    `json:"amount"`
}

type Response struct {
	OrderType     string `json:"orderType"`
	OderID        string `json:"orderID"`
	OrderService  string `json:"orderService"`
	PaymentMethod string `json:"paymentMethod"`
	OrderAmount   int    `json:"orderAmount"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	Amount        int    `json:"amount"`
	ItemId        string `json:"itemId"`
	RespCode      int    `json:"respCode"`
	RespStatus    string `json:"respStatus"`
	RespMessage   string `json:"respMessage"`
}
