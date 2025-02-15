package usecase

import (
	"github.com/nicewook/gocore/internal/domain"
)

type orderUseCase struct {
	orderRepo domain.OrderRepository
}

func NewOrderUseCase(orderRepo domain.OrderRepository) domain.OrderUseCase {
	return &orderUseCase{orderRepo: orderRepo}
}

func (uc *orderUseCase) CreateOrder(order *domain.Order) (*domain.Order, error) {
	return uc.orderRepo.Save(order)
}

func (uc *orderUseCase) GetByID(id int64) (*domain.Order, error) {
	return uc.orderRepo.GetByID(id)
}

func (uc *orderUseCase) GetAll() ([]domain.Order, error) {
	return uc.orderRepo.GetAll()
}
