package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestOrderRepository_Save(t *testing.T) {
	cleanDB(t, "orders", "users", "products")
	ctx := context.Background()

	// Insert user and product to satisfy foreign key constraints
	user := &domain.User{Name: "Test User", Email: "testuser@example.com", Password: "password"}
	product := &domain.Product{Name: "Test Product", PriceInKRW: 1000}

	userRepo := NewUserRepository(testDB)
	productRepo := NewProductRepository(testDB)

	savedUser, err := userRepo.Save(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, savedUser)

	savedProduct, err := productRepo.Save(ctx, product)
	assert.NoError(t, err)
	assert.NotNil(t, savedProduct)

	t.Run("Successfully save order", func(t *testing.T) {
		order := &domain.Order{
			UserID:          savedUser.ID,
			ProductID:       savedProduct.ID,
			Quantity:        2,
			TotalPriceInKRW: 2000,
		}

		orderRepo := NewOrderRepository(testDB)
		savedOrder, err := orderRepo.Save(ctx, order)
		assert.NoError(t, err)
		assert.NotNil(t, savedOrder)
		assert.NotZero(t, savedOrder.ID)
		assert.Equal(t, savedUser.ID, savedOrder.UserID)
		assert.Equal(t, savedProduct.ID, savedOrder.ProductID)
		assert.Equal(t, 2, savedOrder.Quantity)
		assert.Equal(t, int64(2000), savedOrder.TotalPriceInKRW)
		assert.NotEmpty(t, savedOrder.CreatedAt)
	})

	t.Run("Fail to save order with non-existent user", func(t *testing.T) {
		order := &domain.Order{
			UserID:          9999,
			ProductID:       savedProduct.ID,
			Quantity:        2,
			TotalPriceInKRW: 2000,
		}

		orderRepo := NewOrderRepository(testDB)
		savedOrder, err := orderRepo.Save(ctx, order)
		assert.Error(t, err)
		assert.Nil(t, savedOrder)
	})
}

func TestOrderRepository_GetByID(t *testing.T) {
	cleanDB(t, "orders", "users", "products")
	ctx := context.Background()

	// Insert user and product to satisfy foreign key constraints
	user := &domain.User{Name: "Test User", Email: "testuser@example.com", Password: "password"}
	product := &domain.Product{Name: "Test Product", PriceInKRW: 1000}

	userRepo := NewUserRepository(testDB)
	productRepo := NewProductRepository(testDB)

	savedUser, err := userRepo.Save(ctx, user)
	assert.NoError(t, err)

	savedProduct, err := productRepo.Save(ctx, product)
	assert.NoError(t, err)

	// Create an order
	order := &domain.Order{
		UserID:          savedUser.ID,
		ProductID:       savedProduct.ID,
		Quantity:        2,
		TotalPriceInKRW: 2000,
	}

	orderRepo := NewOrderRepository(testDB)
	savedOrder, err := orderRepo.Save(ctx, order)
	assert.NoError(t, err)

	t.Run("Get order by ID successfully", func(t *testing.T) {
		foundOrder, err := orderRepo.GetByID(ctx, savedOrder.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundOrder)
		assert.Equal(t, savedOrder.ID, foundOrder.ID)
		assert.Equal(t, savedUser.ID, foundOrder.UserID)
		assert.Equal(t, savedProduct.ID, foundOrder.ProductID)
		assert.Equal(t, 2, foundOrder.Quantity)
		assert.Equal(t, int64(2000), foundOrder.TotalPriceInKRW)
		assert.NotEmpty(t, foundOrder.CreatedAt)
	})

	t.Run("Order not found", func(t *testing.T) {
		foundOrder, err := orderRepo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.Nil(t, foundOrder)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestOrderRepository_GetAll(t *testing.T) {
	cleanDB(t, "orders", "users", "products")
	ctx := context.Background()

	// Insert user and product to satisfy foreign key constraints
	user := &domain.User{Name: "Test User", Email: "testuser@example.com", Password: "password"}
	product := &domain.Product{Name: "Test Product", PriceInKRW: 1000}

	userRepo := NewUserRepository(testDB)
	productRepo := NewProductRepository(testDB)

	savedUser, err := userRepo.Save(ctx, user)
	assert.NoError(t, err)

	savedProduct, err := productRepo.Save(ctx, product)
	assert.NoError(t, err)

	// Create multiple orders
	orderRepo := NewOrderRepository(testDB)

	order1 := &domain.Order{
		UserID:          savedUser.ID,
		ProductID:       savedProduct.ID,
		Quantity:        1,
		TotalPriceInKRW: 1000,
	}
	savedOrder1, err := orderRepo.Save(ctx, order1)
	assert.NoError(t, err)

	order2 := &domain.Order{
		UserID:          savedUser.ID,
		ProductID:       savedProduct.ID,
		Quantity:        2,
		TotalPriceInKRW: 2000,
	}
	savedOrder2, err := orderRepo.Save(ctx, order2)
	assert.NoError(t, err)

	t.Run("Get all orders successfully", func(t *testing.T) {
		orders, err := orderRepo.GetAll(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, orders)
		assert.Len(t, orders, 2)

		// Check if the returned orders match the saved ones
		foundOrder1 := false
		foundOrder2 := false

		for _, order := range orders {
			if order.ID == savedOrder1.ID {
				foundOrder1 = true
				assert.Equal(t, savedUser.ID, order.UserID)
				assert.Equal(t, savedProduct.ID, order.ProductID)
				assert.Equal(t, 1, order.Quantity)
				assert.Equal(t, int64(1000), order.TotalPriceInKRW)
			}
			if order.ID == savedOrder2.ID {
				foundOrder2 = true
				assert.Equal(t, savedUser.ID, order.UserID)
				assert.Equal(t, savedProduct.ID, order.ProductID)
				assert.Equal(t, 2, order.Quantity)
				assert.Equal(t, int64(2000), order.TotalPriceInKRW)
			}
		}

		assert.True(t, foundOrder1, "First order not found in results")
		assert.True(t, foundOrder2, "Second order not found in results")
	})
}

func TestOrderRepository_GetAll_EmptyResult(t *testing.T) {
	cleanDB(t, "orders")
	ctx := context.Background()
	orderRepo := NewOrderRepository(testDB)

	t.Run("Get all orders with empty result", func(t *testing.T) {
		orders, err := orderRepo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(orders))
	})
}
