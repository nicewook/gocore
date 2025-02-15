package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/config"
	"github.com/nicewook/gocore/internal/db"
	"github.com/nicewook/gocore/internal/handler"
	repository "github.com/nicewook/gocore/internal/repository/postgres"
	"github.com/nicewook/gocore/internal/usecase"
)

func main() {

	env := flag.String("env", "dev", "Environment (dev, qa, stg, prod)")
	flag.Parse()

	validEnvs := map[string]bool{"dev": true, "qa": true, "stg": true, "prod": true}
	if !validEnvs[*env] {
		log.Fatalf("Invalid environment: %s. Valid environments are: dev, qa, stg, prod", *env)
	}

	cfg, err := config.LoadConfig(*env)
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	fmt.Printf("config: %+v\n", cfg)

	// 여기에 DB 연결 및 애플리케이션 로직 추가
	dbConn, err := db.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	// 의존성 주입
	userRepo := repository.NewUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	productRepo := repository.NewProductRepository(dbConn)
	productUseCase := usecase.NewProductUseCase(productRepo)
	productHandler := handler.NewProductHandler(productUseCase)

	// 라우팅
	e := echo.New()

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetByID)
	e.GET("/users", userHandler.GetAll)

	e.POST("/products", productHandler.CreateProduct)
	e.GET("/products/:id", productHandler.GetByID)
	e.GET("/products", productHandler.GetAll)

	// 서버 실행
	log.Println("Server started at :8080")
	if err := e.Start(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
		log.Fatal("Shutting down the server due to:", err)
	}
}
