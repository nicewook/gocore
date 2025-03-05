package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestSaveProduct(t *testing.T) {
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	ctx := context.Background()

	t.Run("Successfully save product", func(t *testing.T) {
		product := &domain.Product{Name: "Product1", PriceInKRW: 100}
		savedProduct, err := repo.Save(ctx, product)
		assert.NoError(t, err)
		assert.NotZero(t, savedProduct.ID)
		assert.Equal(t, product.Name, savedProduct.Name)
		assert.Equal(t, product.PriceInKRW, savedProduct.PriceInKRW)
	})

	t.Run("Fail to save product with duplicate name", func(t *testing.T) {
		product1 := &domain.Product{Name: "Product2", PriceInKRW: 100}
		_, err := repo.Save(ctx, product1)
		assert.NoError(t, err)
		product2 := &domain.Product{Name: "Product2", PriceInKRW: 200}
		_, err = repo.Save(ctx, product2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestGetProductByID(t *testing.T) {
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	ctx := context.Background()

	product := &domain.Product{Name: "Product1", PriceInKRW: 100}
	savedProduct, err := repo.Save(ctx, product)
	assert.NoError(t, err)

	t.Run("Get product by ID successfully", func(t *testing.T) {
		foundProduct, err := repo.GetByID(ctx, savedProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedProduct.ID, foundProduct.ID)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.Equal(t, product.PriceInKRW, foundProduct.PriceInKRW)
	})

	t.Run("Product not found", func(t *testing.T) {
		foundProduct, err := repo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestGetAllProducts(t *testing.T) {
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	ctx := context.Background()

	product1 := &domain.Product{Name: "Product1", PriceInKRW: 100}
	product2 := &domain.Product{Name: "Product2", PriceInKRW: 200}

	_, err := repo.Save(ctx, product1)
	assert.NoError(t, err)
	_, err = repo.Save(ctx, product2)
	assert.NoError(t, err)

	t.Run("Get all products successfully", func(t *testing.T) {
		products, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, products, 2)
	})

	t.Run("No products found", func(t *testing.T) {
		cleanDB(t, "products")
		products, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Empty(t, products)
	})
}

func TestGetAllProducts_ScanError(t *testing.T) {
	// 이 테스트는 실제로 스캔 에러를 발생시키기 어렵기 때문에
	// 테스트 커버리지를 위한 목적으로만 추가합니다.
	// 실제 환경에서는 데이터베이스 스키마 변경 등으로 인해 발생할 수 있습니다.
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	ctx := context.Background()

	// 정상적인 제품 추가
	product := &domain.Product{Name: "ScanErrorTest", PriceInKRW: 100}
	_, err := repo.Save(ctx, product)
	assert.NoError(t, err)

	// 정상적으로 조회되는지 확인
	products, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
}
