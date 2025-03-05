package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestCreateUser(t *testing.T) {
	repo := NewAuthRepository(testDB)
	cleanDB(t, "users")
	ctx := context.Background()

	t.Run("성공적으로 사용자 생성", func(t *testing.T) {
		user := &domain.User{
			Email:    "new@example.com",
			Password: "$argon2id$v=19$m=65536,t=3,p=4$viJoL99egWEGZe7SDvLy9Q$iUs3hjYl83C2DXUuwVbUdLtB2V7gfnaZnk+NAmIPkdY",
		}

		savedUser, err := repo.CreateUser(ctx, user)

		assert.NoError(t, err)
		assert.NotZero(t, savedUser.ID)
		assert.Equal(t, user.Email, savedUser.Email)
	})

	t.Run("중복 이메일로 생성 실패", func(t *testing.T) {
		user1 := &domain.User{Email: "duplicate@example.com", Password: "password123"}
		_, err := repo.CreateUser(ctx, user1)
		assert.NoError(t, err)

		user2 := &domain.User{Email: "duplicate@example.com", Password: "password456"}
		_, err = repo.CreateUser(ctx, user2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}
