package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestSaveUser(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")

	t.Run("성공적으로 사용자 저장", func(t *testing.T) {
		user := &domain.User{Name: "John Doe", Email: "john@example.com"}
		savedUser, err := repo.Save(user)

		assert.NoError(t, err)
		assert.NotZero(t, savedUser.ID)
		assert.Equal(t, user.Name, savedUser.Name)
		assert.Equal(t, user.Email, savedUser.Email)
	})

	t.Run("이메일 중복으로 저장 실패", func(t *testing.T) {
		user1 := &domain.User{Name: "Jane Doe", Email: "jane@example.com"}
		_, err := repo.Save(user1)
		assert.NoError(t, err)

		user2 := &domain.User{Name: "Jane Smith", Email: "jane@example.com"}
		_, err = repo.Save(user2)
		assert.Error(t, err) // UNIQUE 제약 조건 위반
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestGetByID(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")

	t.Run("ID로 사용자 조회 성공", func(t *testing.T) {
		user := &domain.User{Name: "Alice", Email: "alice@example.com"}
		savedUser, _ := repo.Save(user)

		fetchedUser, err := repo.GetByID(savedUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedUser.ID, fetchedUser.ID)
		assert.Equal(t, savedUser.Name, fetchedUser.Name)
		assert.Equal(t, savedUser.Email, fetchedUser.Email)
	})

	t.Run("존재하지 않는 ID로 조회 시 실패", func(t *testing.T) {
		fetchedUser, err := repo.GetByID(9999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, fetchedUser)
	})
}

func TestGetAllUsers(t *testing.T) {
	repo := NewUserRepository(testDB)
	cleanDB(t, "users")

	t.Run("모든 사용자 조회 성공", func(t *testing.T) {
		user1 := &domain.User{Name: "User1", Email: "user1@example.com"}
		user2 := &domain.User{Name: "User2", Email: "user2@example.com"}
		repo.Save(user1)
		repo.Save(user2)

		users, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Contains(t, []string{user1.Name, user2.Name}, users[0].Name)
		assert.Contains(t, []string{user1.Email, user2.Email}, users[0].Email)
		assert.Contains(t, []string{user1.Name, user2.Name}, users[1].Name)
		assert.Contains(t, []string{user1.Email, user2.Email}, users[1].Email)
	})

	t.Run("사용자가 없을 때 빈 배열 반환", func(t *testing.T) {
		cleanDB(t, "users") // 데이터 초기화
		users, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, users, 0)
	})
}
