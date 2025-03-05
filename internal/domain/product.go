package domain

import "context"

type Product struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	PriceInKRW int64  `json:"price_in_krw"`
}

type ProductRepository interface {
	Save(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetAll(ctx context.Context) ([]Product, error)
}

type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetAll(ctx context.Context) ([]Product, error)
}
