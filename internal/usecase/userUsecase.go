package usecase

import (
	"context"

	"github.com/nicewook/gocore/internal/domain"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *userUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.GetAll(ctx)
}
