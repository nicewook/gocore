package domain

import (
	"context"
	"strings"
)

// Role 상수 정의
const (
	RolePublic  = "Public" // 공개 접근 가능, 토큰 불필요
	RoleAdmin   = "Admin"
	RoleManager = "Manager"
	RoleUser    = "User"
)

// GetByIDRequest represents a request to get a user by ID
type GetByIDRequest struct {
	ID int64 `param:"id" validate:"required,min=1"`
}

// GetAllRequest represents a request to get all users with filtering and pagination
type GetAllRequest struct {
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Name   string `query:"name" validate:"omitempty,max=100"`
	Email  string `query:"email" validate:"omitempty,max=100"`
	Roles  string `query:"roles" validate:"omitempty"` // Comma-separated roles
}

// GetRolesArray converts the comma-separated roles string to a slice
func (r *GetAllRequest) GetRolesArray() []string {
	if r.Roles == "" {
		return nil
	}
	return strings.Split(r.Roles, ",")
}

// GetAllResponse represents a paginated response for the GetAll endpoint
type GetAllResponse struct {
	Users      []User `json:"users"`       // 사용자 목록
	TotalCount int64  `json:"total_count"` // 전체 사용자 수 (필터 적용)
	Offset     int    `json:"offset"`      // 현재 오프셋
	Limit      int    `json:"limit"`       // 현재 리밋
	HasMore    bool   `json:"has_more"`    // 더 가져올 데이터가 있는지 여부
}

type User struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name" validate:"omitempty,min=2,max=100"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"-" validate:"required,min=8"`
	Roles    []string `json:"roles"`
}

// RolesToString converts roles slice to comma-separated string for storage
func (u *User) RolesToString() string {
	return strings.Join(u.Roles, ",")
}

// StringToRoles converts comma-separated string to roles slice
func StringToRoles(rolesStr string) []string {
	if rolesStr == "" {
		return []string{}
	}
	return strings.Split(rolesStr, ",")
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdmin checks if user has Admin role
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// IsManager checks if user has Manager role
func (u *User) IsManager() bool {
	return u.HasRole(RoleManager)
}

type UserRepository interface {
	Save(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context, req *GetAllRequest) (*GetAllResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type UserUseCase interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context, req *GetAllRequest) (*GetAllResponse, error)
}
