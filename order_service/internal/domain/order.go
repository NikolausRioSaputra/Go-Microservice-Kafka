package domain

// OrderRequest adalah struktur data yang mewakili request untuk membuat order item
type OrderRequest struct {
	OrderID       string `json:"-"`
	OrderType     string `json:"orderType" binding:"required"`
	TransactionID string `json:"-"`
	UserId        string `json:"userId" binding:"required"`
	ItemId        string `json:"itemId" binding:"required"`
	OrderAmount   int    `json:"orderAmount" binding:"required"`
	PaymentMethod string `json:"paymentMethod"`
	Status        string `json:"-"`
}

// EventRegistrationRequest adalah struktur data yang mewakili request untuk mendaftarkan acara atau layanan.
type EventRegistrationRequest struct {
	OrderID       string `json:"-"`
	EventName     string `json:"eventName" binding:"required"`
	UserID        string `json:"userId" binding:"required"`
	OrderType     string `json:"orderType"`
	TransactionID string `json:"-"`
	Amount        int    `json:"amount" binding:"required"`
	PaymentMethod string `json:"paymentMethod"`
	Status        string `json:"-"`
}
