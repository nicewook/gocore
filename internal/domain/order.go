package domain

import "context"

type Order struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id" validate:"required,gt=0"`
	ProductID       int64  `json:"product_id" validate:"required,gt=0"`
	Quantity        int    `json:"quantity" validate:"required,gt=0"`
	TotalPriceInKRW int64  `json:"total_price_in_krw" validate:"required,gt=0"`
	CreatedAt       string `json:"created_at"`
}

type OrderRepository interface {
	Save(ctx context.Context, order *Order) (*Order, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetAll(ctx context.Context) ([]Order, error)
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetAll(ctx context.Context) ([]Order, error)
}
