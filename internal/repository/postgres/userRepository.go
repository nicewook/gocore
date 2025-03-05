package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/nicewook/gocore/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
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

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, name, email, roles
		FROM users
		WHERE id = $1
	`

	var user domain.User
	var rolesStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &rolesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get user by id")
	}

	user.Roles = domain.StringToRoles(rolesStr)
	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, name, email, roles
		FROM users
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var rolesStr string
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &rolesStr); err != nil {
			return nil, errors.Wrap(err, "failed to scan user")
		}
		user.Roles = domain.StringToRoles(rolesStr)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating users")
	}

	return users, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, name, email, password, roles
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	var rolesStr string
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &rolesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Roles = domain.StringToRoles(rolesStr)
	return user, nil
}
