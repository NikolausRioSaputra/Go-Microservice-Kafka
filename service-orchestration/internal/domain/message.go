package domain

// struct ini di gunakan untuk menangani pesan kafka yang masuk
type IncomingMessage struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService,omitempty"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	ItemId        string `json:"itemId"`
	RespCode      int    `json:"respCode,omitempty"`
	RespStatus    string `json:"respStatus,omitempty"`
	RespMessage   string `json:"respMessage,omitempty"`
}
