package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/nicewook/gocore/internal/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Save(product *domain.Product) (*domain.Product, error) {
	const query = `
		INSERT INTO products (name, price_in_krw)
		VALUES ($1, $2)
		RETURNING id
	`
	if err := r.db.QueryRow(query, product.Name, product.PriceInKRW).Scan(&product.ID); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("product %s: %w", product.Name, domain.ErrAlreadyExists)
		}
		return nil, fmt.Errorf("failed to save product: %w", err)
	}
	return product, nil
}

func (r *productRepository) GetByID(id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, price_in_krw
		FROM products
		WHERE id = $1
	`
	var product domain.Product
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.PriceInKRW)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}
	return &product, nil
}

func (r *productRepository) GetAll() ([]domain.Product, error) {
	query := `
		SELECT id, name, price_in_krw
		FROM products
		ORDER BY id ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()
	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.PriceInKRW); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over products: %w", err)
	}
	return products, nil
}
