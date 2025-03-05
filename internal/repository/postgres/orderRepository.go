package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nicewook/gocore/internal/domain"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	query := `INSERT INTO orders (user_id, product_id, quantity, total_price_in_krw, created_at) 
			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	err := r.DB.QueryRowContext(ctx, query, order.UserID, order.ProductID, order.Quantity, order.TotalPriceInKRW).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `SELECT id, user_id, product_id, quantity, total_price_in_krw, created_at FROM orders WHERE id = $1`
	order := &domain.Order{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.TotalPriceInKRW, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	query := `SELECT id, user_id, product_id, quantity, total_price_in_krw, created_at FROM orders`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.TotalPriceInKRW, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
