package domain

// Struct Message digunakan untuk merepresentasikan data pesan yang diterima.
type Message struct {
	OrderType     string `json:"orderType"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
}

// Struct Response digunakan untuk merepresentasikan data balasan setelah pesan diproses.
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
