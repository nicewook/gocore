package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/nicewook/gocore/internal/config"
	"github.com/nicewook/gocore/internal/db"
	"github.com/nicewook/gocore/internal/handler"
	"github.com/nicewook/gocore/internal/middlewares"
	repository "github.com/nicewook/gocore/internal/repository/postgres"
	"github.com/nicewook/gocore/internal/usecase"
)

func main() {
	app := fx.New(
		fx.Provide(
			NewConfig,
			NewDB,
			echo.New,
		),
		fx.Provide(
			repository.NewUserRepository,
			repository.NewProductRepository,
			repository.NewOrderRepository,
		),
		fx.Provide(
			usecase.NewUserUseCase,
			usecase.NewProductUseCase,
			usecase.NewOrderUseCase,
		),
		fx.Invoke(
			middlewares.RegisterMiddlewares,
			handler.NewUserHandler,
			handler.NewProductHandler,
			handler.NewOrderHandler,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}

func NewConfig() *config.Config {
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
	return cfg
}

func NewDB(lc fx.Lifecycle, cfg *config.Config) *sql.DB {
	dbConn, err := db.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return dbConn.Close()
		},
	})

	return dbConn
}

func StartServer(lc fx.Lifecycle, e *echo.Echo, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
					log.Fatal("Shutting down the server due to:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}
