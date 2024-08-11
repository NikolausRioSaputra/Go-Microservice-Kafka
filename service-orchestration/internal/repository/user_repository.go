package repository

import (
	"context"
	"database/sql"
	"service-orchestration/m/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	// FindByID(ctx context.Context, id string) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := "INSERT INTO users (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Name).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
// 	var user domain.User
// 	query := "SELECT id, name FROM users WHERE id = $1"
// 	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // User not found
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }
