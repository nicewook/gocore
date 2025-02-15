package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestOrderRepository_Save(t *testing.T) {
	cleanDB(t, "orders", "users", "products")

	// Insert user and product to satisfy foreign key constraints
	user := &domain.User{Name: "Test User", Email: "testuser@example.com"}
	product := &domain.Product{Name: "Test Product", PriceInKRW: 1000}

	userRepo := NewUserRepository(testDB)
	productRepo := NewProductRepository(testDB)

	_, err := userRepo.Save(user)
	assert.NoError(t, err)
	_, err = productRepo.Save(product)
	assert.NoError(t, err)

	t.Run("Successfully save order", func(t *testing.T) {
		order := &domain.Order{
			UserID:          1,
			ProductID:       1,
			Quantity:        2,
			TotalPriceInKRW: 2000,
		}

		orderRepo := NewOrderRepository(testDB)
		savedOrder, err := orderRepo.Save(order)

		assert.NoError(t, err)
		assert.NotNil(t, savedOrder)
		assert.Equal(t, order.UserID, savedOrder.UserID)
		assert.Equal(t, order.ProductID, savedOrder.ProductID)
		assert.Equal(t, order.Quantity, savedOrder.Quantity)
		assert.Equal(t, order.TotalPriceInKRW, savedOrder.TotalPriceInKRW)
	})

	t.Run("Fail to save order with non-existent user", func(t *testing.T) {
		order := &domain.Order{
			UserID:          9999,
			ProductID:       1,
			Quantity:        2,
			TotalPriceInKRW: 2000,
		}

		orderRepo := NewOrderRepository(testDB)
		savedOrder, err := orderRepo.Save(order)

		assert.Error(t, err)
		assert.Nil(t, savedOrder)
	})
}

func TestOrderRepository_GetByID(t *testing.T) {
	cleanDB(t, "orders", "users", "products")

	// Insert user and product to satisfy foreign key constraints
	user := &domain.User{Name: "Test User", Email: "testuser@example.com"}
	product := &domain.Product{Name: "Test Product", PriceInKRW: 1000}

	userRepo := NewUserRepository(testDB)
	productRepo := NewProductRepository(testDB)

	_, err := userRepo.Save(user)
	assert.NoError(t, err)
	_, err = productRepo.Save(product)
	assert.NoError(t, err)

	t.Run("Successfully save order", func(t *testing.T) {
		order := &domain.Order{
			UserID:          1,
			ProductID:       1,
			Quantity:        2,
			TotalPriceInKRW: 2000,
		}

		orderRepo := NewOrderRepository(testDB)
		savedOrder, err := orderRepo.Save(order)

		assert.NoError(t, err)
		assert.NotNil(t, savedOrder)
		assert.Equal(t, order.UserID, savedOrder.UserID)
		assert.Equal(t, order.ProductID, savedOrder.ProductID)
		assert.Equal(t, order.Quantity, savedOrder.Quantity)
		assert.Equal(t, order.TotalPriceInKRW, savedOrder.TotalPriceInKRW)
	})

	t.Run("Fail to save order with non-existent user", func(t *testing.T) {
		order := &domain.Order{
			UserID:          9999,
			ProductID:       1,
			Quantity:        2,
			TotalPriceInKRW: 2000,
		}

		orderRepo := NewOrderRepository(testDB)
		savedOrder, err := orderRepo.Save(order)

		assert.Error(t, err)
		assert.Nil(t, savedOrder)
	})
}

func TestOrderRepository_GetAll(t *testing.T) {
	cleanDB(t, "orders", "users", "products")

	user := &domain.User{Name: "Test User", Email: "testuser@example.com"}
	product := &domain.Product{Name: "Test Product", PriceInKRW: 1000}

	// Insert user and product to satisfy foreign key constraints
	_, err := testDB.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, user.Name, user.Email)
	assert.NoError(t, err)
	_, err = testDB.Exec(`INSERT INTO products (name, price_in_krw) VALUES ($1, $2)`, product.Name, product.PriceInKRW)
	assert.NoError(t, err)

	order1 := &domain.Order{
		UserID:          1,
		ProductID:       1,
		Quantity:        2,
		TotalPriceInKRW: 2000,
	}
	order2 := &domain.Order{
		UserID:          1,
		ProductID:       1,
		Quantity:        3,
		TotalPriceInKRW: 3000,
	}

	repo := NewOrderRepository(testDB)
	_, err = repo.Save(order1)
	assert.NoError(t, err)
	_, err = repo.Save(order2)
	assert.NoError(t, err)

	orders, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}
