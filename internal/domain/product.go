package domain

type Product struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	PriceInKRW int64  `json:"price_in_krw"`
}

type ProductRepository interface {
	Save(product *Product) (*Product, error)
	GetByID(id int64) (*Product, error)
	GetAll() ([]Product, error)
}

type ProductUseCase interface {
	CreateProduct(product *Product) (*Product, error)
	GetByID(id int64) (*Product, error)
	GetAll() ([]Product, error)
}
