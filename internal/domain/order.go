package domain

type Order struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	ProductID       int64  `json:"product_id"`
	Quantity        int    `json:"quantity"`
	TotalPriceInKRW int64  `json:"total_price_in_krw"`
	CreatedAt       string `json:"created_at"`
}

type OrderRepository interface {
	Save(order *Order) (*Order, error)
	GetByID(id int64) (*Order, error)
	GetAll() ([]Order, error)
}

type OrderUseCase interface {
	CreateOrder(order *Order) (*Order, error)
	GetByID(id int64) (*Order, error)
	GetAll() ([]Order, error)
}
