package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
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

func (r *userRepository) GetAll(ctx context.Context, req *domain.GetAllRequest) (*domain.GetAllResponse, error) {
	// Squirrel PostgreSQL 쿼리 빌더 초기화
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// 기본 사용자 데이터 쿼리와 카운트 쿼리 생성
	dataBuilder := psql.Select("id", "name", "email", "roles").From("users")
	countBuilder := psql.Select("COUNT(*)").From("users")

	// 필터 조건 적용
	dataBuilder, countBuilder = addUserFilters(dataBuilder, countBuilder, req)

	// 전체 사용자 수 조회
	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "count query 생성 실패")
	}

	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "사용자 수 조회 실패")
	}

	// 결과가 없으면 NotFound 에러 반환
	if totalCount == 0 {
		return nil, domain.ErrNotFound
	}

	// 페이지네이션 및 정렬 적용
	dataBuilder = dataBuilder.
		OrderBy("id ASC").
		Limit(uint64(req.Limit)).
		Offset(uint64(req.Offset))

	// 사용자 데이터 조회
	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "data query 생성 실패")
	}

	rows, err := r.db.QueryContext(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "사용자 데이터 조회 실패")
	}
	defer rows.Close()

	// 결과 처리
	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	// 더 가져올 데이터가 있는지 확인
	hasMore := (req.Offset + len(users)) < int(totalCount)

	return &domain.GetAllResponse{
		Users:      users,
		TotalCount: totalCount,
		Offset:     req.Offset,
		Limit:      req.Limit,
		HasMore:    hasMore,
	}, nil
}

// scanUsers는 DB 결과 세트에서 사용자 목록을 스캔합니다
func scanUsers(rows *sql.Rows) ([]domain.User, error) {
	var users []domain.User

	for rows.Next() {
		var user domain.User
		var rolesStr string

		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &rolesStr); err != nil {
			return nil, errors.Wrap(err, "사용자 스캔 실패")
		}

		user.Roles = domain.StringToRoles(rolesStr)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "행 반복 중 오류 발생")
	}

	return users, nil
}

// addUserFilters는 검색 조건을 쿼리 빌더에 추가합니다
func addUserFilters(dataBuilder, countBuilder squirrel.SelectBuilder, req *domain.GetAllRequest) (squirrel.SelectBuilder, squirrel.SelectBuilder) {
	// 이름 필터
	if req.Name != "" {
		nameFilter := squirrel.ILike{"name": "%" + req.Name + "%"}
		dataBuilder = dataBuilder.Where(nameFilter)
		countBuilder = countBuilder.Where(nameFilter)
	}

	// 이메일 필터
	if req.Email != "" {
		emailFilter := squirrel.ILike{"email": "%" + req.Email + "%"}
		dataBuilder = dataBuilder.Where(emailFilter)
		countBuilder = countBuilder.Where(emailFilter)
	}

	// 역할 필터
	rolesArray := req.GetRolesArray()
	if len(rolesArray) > 0 {
		roleConditions := squirrel.Or{}
		for _, role := range rolesArray {
			roleConditions = append(roleConditions, squirrel.Like{"roles": "%" + role + "%"})
		}

		dataBuilder = dataBuilder.Where(roleConditions)
		countBuilder = countBuilder.Where(roleConditions)
	}

	return dataBuilder, countBuilder
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
