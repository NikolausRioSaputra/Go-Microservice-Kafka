package domain

type Message struct {
	OrderType     string `json:"orderType"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
}

type Response struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
	RespCode      int    `json:"respCode"`
	RespStatus    string `json:"respStatus"`
	RespMessage   string `json:"respMessage"`
}
