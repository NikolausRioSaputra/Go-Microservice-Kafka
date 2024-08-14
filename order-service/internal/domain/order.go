package domain

// OrderRequest adalah struktur data yang mewakili request untuk membuat order baru.
type OrderRequest struct {
	OrderType     string  `json:"orderType" binding:"required"`
	TransactionID string  `json:"transactionId" binding:"required"`
	UserId        string  `json:"userId" binding:"required"`
	PackageId     string  `json:"packageId" binding:"required"`
	OrderAmount   float64 `json:"orderAmount" binding:"required"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
}
