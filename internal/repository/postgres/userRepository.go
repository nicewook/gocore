package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	// 기본 쿼리 설정
	baseQuery := `
		SELECT id, name, email, roles
		FROM users
	`

	// 필터링 조건 추가
	filterQuery, params := r.buildGetAllFilters(req)

	// 결과 갯수를 가져오는 카운트 쿼리
	countQuery := "SELECT COUNT(*) FROM users" + filterQuery

	// 데이터를 가져오는 쿼리 (정렬 및 페이지네이션 추가)
	dataQuery := baseQuery + filterQuery + " ORDER BY id ASC"

	// 페이지네이션 추가
	dataQuery += " LIMIT $" + fmt.Sprintf("%d", len(params)+1) + " OFFSET $" + fmt.Sprintf("%d", len(params)+2)
	params = append(params, req.Limit, req.Offset)

	// 전체 사용자 수 조회
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, params[:len(params)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total user count")
	}

	// 데이터 조회
	rows, err := r.db.QueryContext(ctx, dataQuery, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
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

	// 결과가 없는 경우 처리
	if len(users) == 0 && totalCount == 0 {
		return nil, domain.ErrNotFound
	}

	// 더 가져올 데이터가 있는지 확인
	hasMore := (req.Offset + len(users)) < int(totalCount)

	// 응답 생성
	response := &domain.GetAllResponse{
		Users:      users,
		TotalCount: totalCount,
		Offset:     req.Offset,
		Limit:      req.Limit,
		HasMore:    hasMore,
	}

	return response, nil
}

// buildGetAllFilters는 사용자 필터링을 위한 SQL 조건절과 파라미터를 생성합니다
func (r *userRepository) buildGetAllFilters(req *domain.GetAllRequest) (string, []interface{}) {
	var conditions []string
	var params []interface{}
	paramCount := 1

	// WHERE 절 구성
	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", paramCount))
		params = append(params, "%"+req.Name+"%")
		paramCount++
	}

	if req.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", paramCount))
		params = append(params, "%"+req.Email+"%")
		paramCount++
	}

	if req.Roles != "" {
		// roles 필드는 콤마로 구분된 문자열로 저장되어 있으므로
		// 각 역할에 대해 LIKE 조건을 생성
		rolesArray := req.GetRolesArray()
		if len(rolesArray) > 0 {
			var roleConditions []string
			for _, role := range rolesArray {
				roleConditions = append(roleConditions, fmt.Sprintf("roles LIKE $%d", paramCount))
				params = append(params, "%"+role+"%")
				paramCount++
			}
			if len(roleConditions) > 0 {
				conditions = append(conditions, "("+strings.Join(roleConditions, " OR ")+")")
			}
		}
	}

	// WHERE 절 생성
	var filterQuery string
	if len(conditions) > 0 {
		filterQuery = " WHERE " + strings.Join(conditions, " AND ")
	}

	return filterQuery, params
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
