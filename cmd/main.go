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
	// DI
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	otpRepo := repository.NewOTPRepository(config.DB)
	mailService := service.NewMailService(otpRepo)
	authHandler := handler.NewAuthHandler(userService, mailService)

	e := echo.New()
	routes.RegisterRoutes(e, userHandler, authHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
