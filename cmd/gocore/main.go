package main

import (
	"log"

	"github.com/nicewook/gocore/internal/handler"
	"github.com/nicewook/gocore/internal/repository"
	"github.com/nicewook/gocore/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// 의존성 주입
	userRepo := repository.NewUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	// 라우팅
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)

	// 서버 실행
	log.Println("Server started at :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Shutting down the server due to:", err)
	}
}
