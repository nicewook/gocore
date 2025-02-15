package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nicewook/gocore/internal/config"
)

func NewDBConnection(cfg *config.Config) (*sql.DB, error) {
	// DSN(Data Source Name) 생성
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode,
	)

	// 데이터베이스 연결 생성: sql.Open은 실제 연결을 생성하는 것이 아니라 연결 가능한 객체를 반환한다. db.Ping()을 통해 실제 연결을 확인한다.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 연결 풀 설정 (성능 최적화)
	db.SetMaxOpenConns(25)                 // 동시에 열 수 있는 최대 연결 수
	db.SetMaxIdleConns(25)                 // 유휴 상태로 유지할 연결의 최대 수
	db.SetConnMaxLifetime(5 * time.Minute) // 연결의 최대 수명 (5분 후 재연결)

	// 실제 연결을 테스트하여 연결 가능 여부 확인
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"failed to ping database (host: %s, db: %s): %w",
			cfg.DB.Host, cfg.DB.DBName, err,
		)
	}

	if err := createUserTable(db); err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	if err := createProductTable(db); err != nil {
		return nil, fmt.Errorf("failed to create products table: %w", err)
	}

	if err := createOrderTable(db); err != nil {
		return nil, fmt.Errorf("failed to create orders table: %w", err)
	}

	return db, nil
}

// DB 생성시에 User 테이블을 생성하는 함수(존재하지 않을 경우에는)
func createUserTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE
		)
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

func createProductTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
            price_in_krw BIGINT NOT NULL,
			UNIQUE (name)
		)
	`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}
	return nil
}

func createOrderTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			product_id INT NOT NULL,
			quantity INT NOT NULL,
		    total_price_in_krw BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)
	`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}
	return nil
}
