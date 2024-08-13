package domain

// Struct Message digunakan untuk merepresentasikan data pesan yang diterima.

// Struct Message digunakan untuk merepresentasikan data pesan yang diterima.
type PaymentMessage struct {
	OrderType     string  `json:"orderType"`
	TransactionId string  `json:"transactionId"`
	UserId        string  `json:"userId"`
	PackageId     string  `json:"packageId"`
	OrderAmount   float64 `json:"orderAmount"`
	PaymentMethod string  `json:"paymentMethod"`
}

type PaymentResponse struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
	RespCode      int    `json:"respCode"`
	RespStatus    string `json:"respStatus"`
	RespMessage   string `json:"respMessage"`
}
