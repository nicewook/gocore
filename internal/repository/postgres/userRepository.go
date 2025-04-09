package postgres

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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

func likeCondition(key string, value interface{}) sq.Like {
	if value == nil || value == "" {
		return nil
	}
	return sq.Like{
		key: fmt.Sprintf("%%%s%%", value),
	}
}

func equalCondition(key string, value interface{}) sq.Eq {
	if value == nil {
		return nil
	}
	return sq.Eq{
		key: value,
	}
}

func (r *userRepository) GetAll(ctx context.Context, req *domain.GetAllUsersRequest) (*domain.GetAllResponse, error) {

	// PostgreSQL 스타일의 파라미터 사용을 위한 placeholder 설정
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// 사용자 데이터 쿼리, 카운트 쿼리 빌더 생성
	dataBuilder := psql.Select("id", "name", "email", "password", "roles").From("users")
	countBuilder := psql.Select("COUNT(*)").From("users")

	// 필터 조건 적용
	if req.Name != "" {
		dataBuilder = dataBuilder.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", req.Name)})
		countBuilder = countBuilder.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", req.Name)})
	}
	if req.Email != "" {
		dataBuilder = dataBuilder.Where(sq.Like{"email": fmt.Sprintf("%%%s%%", req.Email)})
		countBuilder = countBuilder.Where(sq.Like{"email": fmt.Sprintf("%%%s%%", req.Email)})
	}

	rolesArray := req.GetRolesArray()
	if len(rolesArray) > 0 {
		roleConditions := make([]sq.Sqlizer, 0, len(rolesArray))
		for _, role := range rolesArray {
			if role != "" {
				roleConditions = append(roleConditions, sq.Like{"roles": fmt.Sprintf("%%%s%%", role)})
			}
		}
		if len(roleConditions) > 0 {
			dataBuilder = dataBuilder.Where(sq.Or(roleConditions))
			countBuilder = countBuilder.Where(sq.Or(roleConditions))
		}
	}

	// 카운트 쿼리 실행
	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "count query 생성 실패")
	}

	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "사용자 수 조회 실패")
	}

	if totalCount == 0 {
		// 사용자가 없을 때는 빈 배열과 함께 응답 반환
		return &domain.GetAllResponse{
			Users:      []domain.User{},
			TotalCount: 0,
			Offset:     req.Offset,
			Limit:      req.Limit,
			HasMore:    false,
		}, nil
	}

	// 사용자 데이터 쿼리 페이지네이션 및 정렬 적용
	dataBuilder = dataBuilder.
		OrderBy("id ASC").
		Limit(uint64(req.Limit)).
		Offset(uint64(req.Offset))

	// 사용자 데이터 쿼리 실행
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
		var password string

		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &password, &rolesStr); err != nil {
			return nil, errors.Wrap(err, "사용자 스캔 실패")
		}

		user.Password = password
		user.Roles = domain.StringToRoles(rolesStr)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "행 반복 중 오류 발생")
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
