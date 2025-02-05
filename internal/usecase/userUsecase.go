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

func (uc *userUseCase) CreateUser(user domain.User) error {
	if user.Name == "" || user.Email == "" {
		return domain.ErrInvalidInput
	}
	return uc.userRepo.Save(user)
}

func (uc *userUseCase) GetUser(id int) (domain.User, error) {
	return uc.userRepo.FindByID(id)
}
