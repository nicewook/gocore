package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
)

func TestSaveProduct(t *testing.T) {
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	t.Run("Successfully save product", func(t *testing.T) {
		product := &domain.Product{Name: "Product1", PriceInKRW: 100}
		savedProduct, err := repo.Save(product)
		assert.NoError(t, err)
		assert.NotZero(t, savedProduct.ID)
		assert.Equal(t, product.Name, savedProduct.Name)
		assert.Equal(t, product.PriceInKRW, savedProduct.PriceInKRW)
	})
	t.Run("Fail to save product with duplicate name", func(t *testing.T) {
		product1 := &domain.Product{Name: "Product2", PriceInKRW: 100}
		_, err := repo.Save(product1)
		assert.NoError(t, err)
		product2 := &domain.Product{Name: "Product2", PriceInKRW: 200}
		_, err = repo.Save(product2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestGetProductByID(t *testing.T) {
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	t.Run("Successfully get product by ID", func(t *testing.T) {
		product := &domain.Product{Name: "Product1", PriceInKRW: 100}
		savedProduct, _ := repo.Save(product)
		fetchedProduct, err := repo.GetByID(savedProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedProduct.ID, fetchedProduct.ID)
		assert.Equal(t, savedProduct.Name, fetchedProduct.Name)
		assert.Equal(t, savedProduct.PriceInKRW, fetchedProduct.PriceInKRW)
	})
	t.Run("Fail to get product by non-existent ID", func(t *testing.T) {
		fetchedProduct, err := repo.GetByID(9999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, fetchedProduct)
	})
}

func TestGetAllProducts(t *testing.T) {
	repo := NewProductRepository(testDB)
	cleanDB(t, "products")
	t.Run("Successfully get all products", func(t *testing.T) {
		product1 := &domain.Product{Name: "Product1", PriceInKRW: 100}
		product2 := &domain.Product{Name: "Product2", PriceInKRW: 200}
		repo.Save(product1)
		repo.Save(product2)
		products, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Contains(t, []string{product1.Name, product2.Name}, products[0].Name)
		assert.Contains(t, []int64{product1.PriceInKRW, product2.PriceInKRW}, products[0].PriceInKRW)
		assert.Contains(t, []string{product1.Name, product2.Name}, products[1].Name)
		assert.Contains(t, []int64{product1.PriceInKRW, product2.PriceInKRW}, products[1].PriceInKRW)
	})
	t.Run("Return empty array when no products found", func(t *testing.T) {
		cleanDB(t, "products")
		products, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, products, 0)
	})
}
