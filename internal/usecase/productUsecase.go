package usecase

import (
	"context"

	"github.com/nicewook/gocore/internal/domain"
)

type productUseCase struct {
	productRepo domain.ProductRepository
}

func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{productRepo: productRepo}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product.Name == "" || product.PriceInKRW <= 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.productRepo.Save(ctx, product)
}

func (uc *productUseCase) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}

func (uc *productUseCase) GetAll(ctx context.Context) ([]domain.Product, error) {
	return uc.productRepo.GetAll(ctx)
}
