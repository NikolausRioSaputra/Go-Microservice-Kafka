package domain

// OrderRequest adalah struktur data yang mewakili request untuk membuat order baru.
type OrderRequest struct {
	OrderType     string `json:"orderType" binding:"required"`
	TransactionID string `json:"-"`
	UserId        string `json:"userId" binding:"required"`
	ItemId        string `json:"itemId" binding:"required"`
	OrderAmount   int    `json:"orderAmount" binding:"required"`
	PaymentMethod string `json:"paymentMethod" binding:"required"`
}
