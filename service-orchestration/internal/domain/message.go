package domain
// struct ini di gunakan untuk menangani pesan kafka yang masuk
type IncomingMessage struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService,omitempty"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
	RespCode      int    `json:"respCode,omitempty"`
	RespStatus    string `json:"respStatus,omitempty"`
	RespMessage   string `json:"respMessage,omitempty"`
}
