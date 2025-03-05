package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/nicewook/gocore/internal/domain"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// 기본 역할이 없으면 User 역할 추가
	if len(user.Roles) == 0 {
		user.Roles = []string{domain.RoleUser}
	}

	const query = `
		INSERT INTO users (name, email, password, roles)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	if err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.RolesToString()).Scan(&user.ID); err != nil {
		// PostgreSQL의 unique_violation 에러 코드 (23505)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("email %s: %w", user.Email, domain.ErrAlreadyExists)
		}
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}
