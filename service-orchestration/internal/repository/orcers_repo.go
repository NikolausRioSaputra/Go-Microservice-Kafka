package repository

import (
	"database/sql"
)

type OcresRepository interface {
	ViewTopic(orderType, orderService string) (string, error)
}

type ocresRepository struct {
	DB *sql.DB
}

func NewViewTopicRepository(db *sql.DB) OcresRepository {
	return &ocresRepository{
		DB: db,
	}
}

// viewTopic implements OrderRepository.
func (repo *ocresRepository) ViewTopic(orderType string, orderService string) (string ,error) {
	var topic string
	query := `select topic
			from orces_step
			where order_type = $1 and order_service = $2`
	err := repo.DB.QueryRow(query, orderType, orderService).Scan(&topic)
	if err != nil {
		return "", err
	}
	return topic, nil
}
