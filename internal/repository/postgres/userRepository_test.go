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

		users, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Contains(t, []string{user1.Name, user2.Name}, users[0].Name)
		assert.Contains(t, []string{user1.Email, user2.Email}, users[0].Email)
		assert.Contains(t, []string{user1.Name, user2.Name}, users[1].Name)
		assert.Contains(t, []string{user1.Email, user2.Email}, users[1].Email)
	})

	t.Run("사용자가 없을 때 빈 배열 반환", func(t *testing.T) {
		cleanDB(t, "users") // 데이터 초기화
		users, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, users, 0)
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
