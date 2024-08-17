package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"service-orchestration/m/internal/domain"
)

type OcresRepository interface {
	ViewTopic(orderType, orderService string) (string, error)
	SaveTransaction(message domain.Message, topic, stepStatus string) (int, error)
}

type ocresRepository struct {
	DB *sql.DB
}

func NewOcresRepository(db *sql.DB) OcresRepository {
	return &ocresRepository{
		DB: db,
	}
}

// viewTopic implements OrderRepository.
func (repo *ocresRepository) ViewTopic(orderType string, orderService string) (string, error) {
	var topic string
	query := `SELECT topic
			FROM orces_step
			WHERE order_type = $1 AND order_service = $2`
	err := repo.DB.QueryRow(query, orderType, orderService).Scan(&topic)
	if err != nil {
		return "", err
	}
	return topic, nil
}

func (repo *ocresRepository) SaveTransaction(message domain.Message, topic, stepStatus string) (int, error) {
	var id int
	query := `
		INSERT INTO t_transactions (
			transaction_id, order_id, order_type, order_service, topic, step_status,
			balance, payment_method, order_amount, price, user_id, item_id,
			resp_code, resp_status, resp_message, payload
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id`

	// Serialize the entire message to JSON for the payload column
	payload, err := json.Marshal(message)
	if err != nil {
		return 0, fmt.Errorf("error serializing payload: %v", err)
	}

	err = repo.DB.QueryRow(
		query,
		message.TransactionId,
		message.OderID,
		message.OrderType,
		message.OrderService,
		topic,
		stepStatus,
		message.Balance,
		message.PaymentMethod,
		message.OrderAmount,
		message.Price,
		message.UserId,
		message.ItemId,
		message.RespCode,
		message.RespStatus,
		message.RespMessage,
		payload,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error saving transaction: %v", err)
	}

	return id, nil
}
