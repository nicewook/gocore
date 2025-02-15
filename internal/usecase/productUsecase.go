package usecase

import (
	"github.com/nicewook/gocore/internal/domain"
)

type productUseCase struct {
	productRepo domain.ProductRepository
}

func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{productRepo: productRepo}
}

func (uc *productUseCase) CreateProduct(product *domain.Product) (*domain.Product, error) {
	if product.Name == "" || product.PriceInKRW <= 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.productRepo.Save(product)
}

func (uc *productUseCase) GetByID(id int64) (*domain.Product, error) {
	return uc.productRepo.GetByID(id)
}

func (uc *productUseCase) GetAll() ([]domain.Product, error) {
	return uc.productRepo.GetAll()
}
