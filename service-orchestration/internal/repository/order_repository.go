package repository

import (
	"context"
	"database/sql"
	"service-orchestration/m/internal/domain"
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
	query := `INSERT INTO orders (order_type, transaction_id, user_id, package_id) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.ExecContext(ctx, query, order.OrderType, order.TransactionID, order.UserId, order.PackageId)
	return err
}
