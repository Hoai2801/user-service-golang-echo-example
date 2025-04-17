package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"user-service/config"
	"user-service/internal/handler"
	"user-service/middleware"
)

func RegisterRoutes(e *echo.Echo, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	v1 := e.Group("/api/v1")

	userGroup := v1.Group("/users")
	userGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: config.GetJWTSecret(),
	}))
	userGroup.Use(middleware.AuthMiddleware)
	userGroup.GET("", userHandler.GetAllUsers)
	userGroup.POST("", userHandler.CreateUser)

	authGroup := v1.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/logout", authHandler.Logout)
}
