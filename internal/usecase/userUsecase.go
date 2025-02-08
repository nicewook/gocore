package usecase

import (
	"github.com/nicewook/gocore/internal/domain"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) CreateUser(user *domain.User) (*domain.User, error) {
	if user.Name == "" || user.Email == "" {
		return nil, domain.ErrInvalidInput
	}
	return uc.userRepo.Save(user)
}

func (uc *userUseCase) GetByID(id int64) (*domain.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *userUseCase) GetAll() ([]domain.User, error) {
	return uc.userRepo.GetAll()
}
