package usecase

import (
	"context"
	"service-orchestration/m/internal/domain"
	"service-orchestration/m/internal/repository"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (uc *userUseCase) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return uc.userRepo.Save(ctx, user)
}

// func (uc *userUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
// 	return uc.userRepo.FindByID(ctx, id)
// }
