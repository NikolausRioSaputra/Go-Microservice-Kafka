package repository

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order domain.OrderRequest) (domain.OrderRequest, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
}

type orderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) SaveOrder(ctx context.Context, order domain.OrderRequest) (domain.OrderRequest, error) {
	query := `INSERT INTO orders (order_type, transaction_id, user_id, item_id, order_amount, payment_method, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id`

	err := r.DB.QueryRowContext(ctx, query,
		order.OrderType,
		order.TransactionID,
		order.UserId,
		order.ItemId,
		order.OrderAmount,
		order.PaymentMethod,
		"created").Scan(&order.OrderID) // Set initial status as "pending"

	if err != nil {
		return domain.OrderRequest{}, err
	}

	order.Status = "created"
	return order, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := r.DB.ExecContext(ctx, query, status, orderID)
	return err
}
