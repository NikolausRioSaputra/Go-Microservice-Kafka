package repository

import (
	"database/sql"
	"fmt"
)

type OcresRepository interface {
	ViewTopic(orderType, orderService string) (string, error)
	SaveTransaction(transactionID, orderType, orderService, topic, stepStatus string) (int, error)
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

func (repo *ocresRepository) SaveTransaction(transactionID, orderType, orderService, topic, stepStatus string) (int, error) {
	var id int
	query := `
		INSERT INTO t_transactions (transaction_id, order_type, order_service, topic, step_status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	// Menggunakan QueryRow untuk menyimpan data dan mengambil id yang baru saja di-generate
	err := repo.DB.QueryRow(query, transactionID, orderType, orderService, topic, stepStatus).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error saving transaction: %v", err)
	}

	return id, nil
}
