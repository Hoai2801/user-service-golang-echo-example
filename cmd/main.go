package main

import (
	"github.com/labstack/echo/v4"
	"user-service/config"
	"user-service/internal/cache"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/routes"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()
	config.InitDB()
	cache.InitializeRedis()
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authHandler := handler.NewAuthHandler(userService)

	e := echo.New()
	routes.RegisterRoutes(e, userHandler, authHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
