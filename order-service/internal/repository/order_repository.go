package repository

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order domain.OrderRequest) error
}

type orderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) SaveOrder(ctx context.Context, order domain.OrderRequest) error {
	query := `INSERT INTO orders (order_type, transaction_id, user_id, item_id, order_amount, payment_method)
              VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.DB.ExecContext(ctx, query,
		order.OrderType,
		order.TransactionID,
		order.UserId,
		order.ItemId,
		order.OrderAmount,   // Pastikan ada OrderAmount dalam domain.OrderRequest
		order.PaymentMethod) // Pastikan ada PaymentMethod dalam domain.OrderRequest

	return err
}
