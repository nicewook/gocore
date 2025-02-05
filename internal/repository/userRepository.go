package repository

import (
	"github.com/nicewook/gocore/internal/domain"
)

type userRepository struct {
	users map[int]domain.User // 간단한 메모리 저장소
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{
		users: make(map[int]domain.User),
	}
}

func (r *userRepository) Save(user domain.User) error {
	if _, exists := r.users[user.ID]; exists {
		return domain.ErrAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) FindByID(id int) (domain.User, error) {
	user, exists := r.users[id]
	if !exists {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}
