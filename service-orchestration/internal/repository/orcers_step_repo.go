package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"service-orchestration/m/internal/domain"
)

type OcresRepository interface {
	ViewTopic(orderType, orderService string) (string, error)
	SaveTransaction(message domain.Message, topic, stepStatus string) (int, error)
	GetAllTransactions() ([]domain.Message, error)
	UpdateTransactionItemId(transactionId, newItemId string) error
	GetTransactionByID(transactionId string) (*domain.Message, error)
	UpdateTransactionPayload(transactionId, updatedPayload, newOrderService string) error
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

func (repo *ocresRepository) GetAllTransactions() ([]domain.Message, error) {
	query := `SELECT transaction_id, order_id, order_type, order_service, balance, payment_method, order_amount, price, user_id, item_id, resp_code, resp_status, resp_message , payload FROM t_transactions`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Message
	for rows.Next() {
		var transaction domain.Message
		err := rows.Scan(
			&transaction.TransactionId,
			&transaction.OderID,
			&transaction.OrderType,
			&transaction.OrderService,
			&transaction.Balance,
			&transaction.PaymentMethod,
			&transaction.OrderAmount,
			&transaction.Price,
			&transaction.UserId,
			&transaction.ItemId,
			&transaction.RespCode,
			&transaction.RespStatus,
			&transaction.RespMessage,
			&transaction.Payload,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repo *ocresRepository) UpdateTransactionItemId(transactionId, newItemId string) error {
	query := `UPDATE t_transactions SET item_id = $1 WHERE transaction_id = $2`

	_, err := repo.DB.Exec(query, newItemId, transactionId)
	if err != nil {
		return fmt.Errorf("error updating transaction itemId: %v", err)
	}

	return nil
}

func (repo *ocresRepository) UpdateTransactionPayload(transactionId, updatedPayload, newOrderService string) error {
    query := `UPDATE t_transactions SET payload = $1, order_service = $2 WHERE transaction_id = $3`

    _, err := repo.DB.Exec(query, updatedPayload, newOrderService, transactionId)
    if err != nil {
        return fmt.Errorf("error updating transaction payload: %v", err)
    }

    return nil
}

func (repo *ocresRepository) GetTransactionByID(transactionId string) (*domain.Message, error) {
    var transaction domain.Message
    query := `SELECT transaction_id, order_id, order_type, order_service, balance, payment_method, order_amount, price, user_id, item_id, resp_code, resp_status, resp_message, payload
              FROM t_transactions
              WHERE transaction_id = $1`
    err := repo.DB.QueryRow(query, transactionId).Scan(
        &transaction.TransactionId,
        &transaction.OderID,
        &transaction.OrderType,
        &transaction.OrderService,
        &transaction.Balance,
        &transaction.PaymentMethod,
        &transaction.OrderAmount,
        &transaction.Price,
        &transaction.UserId,
        &transaction.ItemId,
        &transaction.RespCode,
        &transaction.RespStatus,
        &transaction.RespMessage,
        &transaction.Payload,
    )
    if err != nil {
        return nil, fmt.Errorf("error getting transaction: %v", err)
    }

    // Log the retrieved transaction for debugging
    log.Printf("Retrieved transaction: %+v", transaction)

    return &transaction, nil
}
