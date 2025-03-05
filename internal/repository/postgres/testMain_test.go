package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB

func TestMain(m *testing.M) {

	// PostgreSQL 컨테이너 시작
	var (
		ctx       = context.Background()
		container testcontainers.Container
		err       error
	)
	container, testDB, err = setupPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("PostgreSQL 컨테이너 설정 실패: %v", err)
	}

	setupSchema()   // 스키마 생성
	code := m.Run() // 테스트 실행

	// 테스트 종료 후 컨테이너 정리
	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("컨테이너 종료 실패: %v", err)
	}

	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.2", // 원하는 PostgreSQL 버전
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("컨테이너 시작 실패: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, err
	}

	// DB 연결이 준비될 때까지 대기
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return container, db, nil
}

func setupSchema() {
	const schema = `
        CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			roles VARCHAR(255) DEFAULT 'User'
		);
        CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
            price_in_krw BIGINT NOT NULL,
			UNIQUE (name)
		);
        CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			product_id INT NOT NULL,
			quantity INT NOT NULL,
		    total_price_in_krw BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		);
    `
	if _, err := testDB.Exec(schema); err != nil {
		log.Fatalf("테이블 생성 실패: %v", err)
	}
}

// 데이터 초기화 함수
func cleanDB(t *testing.T, tables ...string) {
	for _, table := range tables {
		// "TRUNCATE TABLE " + table 은 인텔리제이가 에러라고 생각한다.
		_, err := testDB.Exec("TRUNCATE TABLE" + " " + table + " RESTART IDENTITY CASCADE")
		assert.NoError(t, err)
	}
}
