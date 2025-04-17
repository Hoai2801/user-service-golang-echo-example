package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"user-service/config"
	"user-service/internal/cache"
	"user-service/internal/model"
)

// AuthMiddleware is the middleware for checking JWT authentication
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing Authorization Header"})
		}

		// Extract the token from the header
		tokenString := authHeader[len("Bearer "):]

		// Check if token is blacklisted
		isBlacklisted, err := cache.CheckBlacklisted(tokenString)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not check blacklist"})
		}
		if isBlacklisted {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token is blacklisted"})
		}

		// Check if token is valid and cached
		isValid, err := cache.CheckTokenInCache(tokenString)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not check token cache"})
		}
		if !isValid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})
		fmt.Println(token)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
		}

		// Set user info in context (optional)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", &model.User{
				Name:  claims["username"].(string),
				Email: claims["email"].(string),
			})
		}

		// Proceed to next handler
		return next(c)
	}
}
