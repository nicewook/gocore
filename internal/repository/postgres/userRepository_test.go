package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestSaveUser(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")
	ctx := context.Background()

	t.Run("성공적으로 사용자 저장", func(t *testing.T) {
		user := &domain.User{Name: "John Doe", Email: "john@example.com"}
		savedUser, err := repo.Save(ctx, user)

		assert.NoError(t, err)
		assert.NotZero(t, savedUser.ID)
		assert.Equal(t, user.Name, savedUser.Name)
		assert.Equal(t, user.Email, savedUser.Email)
	})

	t.Run("이메일 중복으로 저장 실패", func(t *testing.T) {
		user1 := &domain.User{Name: "Jane Doe", Email: "jane@example.com"}
		_, err := repo.Save(ctx, user1)
		assert.NoError(t, err)

		user2 := &domain.User{Name: "Jane Smith", Email: "jane@example.com"}
		_, err = repo.Save(ctx, user2)
		assert.Error(t, err) // UNIQUE 제약 조건 위반
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestGetByID(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")
	ctx := context.Background()

	t.Run("ID로 사용자 조회 성공", func(t *testing.T) {
		user := &domain.User{Name: "Alice", Email: "alice@example.com"}
		savedUser, _ := repo.Save(ctx, user)

		fetchedUser, err := repo.GetByID(ctx, savedUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedUser.ID, fetchedUser.ID)
		assert.Equal(t, savedUser.Name, fetchedUser.Name)
		assert.Equal(t, savedUser.Email, fetchedUser.Email)
	})

	t.Run("존재하지 않는 ID로 조회 시 실패", func(t *testing.T) {
		fetchedUser, err := repo.GetByID(ctx, 9999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, fetchedUser)
	})
}

func TestGetAllUsers(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")
	ctx := context.Background()

	t.Run("모든 사용자 조회 성공", func(t *testing.T) {
		user1 := &domain.User{Name: "User1", Email: "user1@example.com"}
		user2 := &domain.User{Name: "User2", Email: "user2@example.com"}
		repo.Save(ctx, user1)
		repo.Save(ctx, user2)

		// 빈 요청으로 모든 사용자 조회
		req := &domain.GetAllRequest{
			Limit: 10,
		}
		response, err := repo.GetAll(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Users, 2)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, 0, response.Offset)
		assert.Equal(t, 10, response.Limit)
		assert.False(t, response.HasMore)

		// 사용자 데이터 확인
		userNames := []string{response.Users[0].Name, response.Users[1].Name}
		userEmails := []string{response.Users[0].Email, response.Users[1].Email}
		assert.Contains(t, userNames, "User1")
		assert.Contains(t, userNames, "User2")
		assert.Contains(t, userEmails, "user1@example.com")
		assert.Contains(t, userEmails, "user2@example.com")
	})

	t.Run("필터링 기능 테스트", func(t *testing.T) {
		cleanDB(t, "users") // 데이터 초기화

		// 테스트 데이터 생성
		user1 := &domain.User{Name: "John Doe", Email: "john@example.com", Roles: []string{domain.RoleUser}}
		user2 := &domain.User{Name: "Jane Smith", Email: "jane@example.com", Roles: []string{domain.RoleManager}}
		user3 := &domain.User{Name: "Admin User", Email: "admin@example.com", Roles: []string{domain.RoleAdmin}}

		repo.Save(ctx, user1)
		repo.Save(ctx, user2)
		repo.Save(ctx, user3)

		// 이름으로 필터링
		nameReq := &domain.GetAllRequest{
			Name:  "John",
			Limit: 10,
		}
		nameResponse, err := repo.GetAll(ctx, nameReq)
		assert.NoError(t, err)
		assert.Len(t, nameResponse.Users, 1)
		assert.Equal(t, "John Doe", nameResponse.Users[0].Name)

		// 이메일로 필터링
		emailReq := &domain.GetAllRequest{
			Email: "admin",
			Limit: 10,
		}
		emailResponse, err := repo.GetAll(ctx, emailReq)
		assert.NoError(t, err)
		assert.Len(t, emailResponse.Users, 1)
		assert.Equal(t, "admin@example.com", emailResponse.Users[0].Email)

		// 페이지네이션 테스트
		pageReq := &domain.GetAllRequest{
			Offset: 0,
			Limit:  2,
		}
		pageResponse, err := repo.GetAll(ctx, pageReq)
		assert.NoError(t, err)
		assert.Len(t, pageResponse.Users, 2)
		assert.Equal(t, int64(3), pageResponse.TotalCount)
		assert.True(t, pageResponse.HasMore)

		// 두 번째 페이지
		page2Req := &domain.GetAllRequest{
			Offset: 2,
			Limit:  2,
		}
		page2Response, err := repo.GetAll(ctx, page2Req)
		assert.NoError(t, err)
		assert.Len(t, page2Response.Users, 1)
		assert.Equal(t, int64(3), page2Response.TotalCount)
		assert.False(t, page2Response.HasMore)
	})

	t.Run("사용자가 없을 때 에러 반환", func(t *testing.T) {
		cleanDB(t, "users") // 데이터 초기화

		req := &domain.GetAllRequest{
			Limit: 10,
		}
		response, err := repo.GetAll(ctx, req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, response)
	})
}

func TestGetUserByEmail(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")
	ctx := context.Background()

	t.Run("이메일로 사용자 조회 성공", func(t *testing.T) {
		// 테스트 사용자 생성
		user := &domain.User{
			Name:     "Email Test User",
			Email:    "emailtest@example.com",
			Password: "hashedpassword123",
		}
		savedUser, err := repo.Save(ctx, user)
		assert.NoError(t, err)

		// 이메일로 사용자 조회
		fetchedUser, err := repo.GetUserByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, fetchedUser)
		assert.Equal(t, savedUser.ID, fetchedUser.ID)
		assert.Equal(t, user.Name, fetchedUser.Name)
		assert.Equal(t, user.Email, fetchedUser.Email)
		assert.Equal(t, "hashedpassword123", fetchedUser.Password)
		assert.Contains(t, fetchedUser.Roles, domain.RoleUser)
	})

	t.Run("존재하지 않는 이메일로 조회 시 실패", func(t *testing.T) {
		// 존재하지 않는 이메일로 조회
		fetchedUser, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, fetchedUser)
	})
}
