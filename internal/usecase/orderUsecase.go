package usecase

import (
	"context"

	"github.com/nicewook/gocore/internal/domain"
)

type orderUseCase struct {
	orderRepo domain.OrderRepository
}

func NewOrderUseCase(orderRepo domain.OrderRepository) domain.OrderUseCase {
	return &orderUseCase{orderRepo: orderRepo}
}

func (uc *orderUseCase) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	return uc.orderRepo.Save(ctx, order)
}

func (uc *orderUseCase) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	return uc.orderRepo.GetByID(ctx, id)
}

func (uc *orderUseCase) GetAll(ctx context.Context) ([]domain.Order, error) {
	return uc.orderRepo.GetAll(ctx)
}
