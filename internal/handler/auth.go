package handler

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
	"user-service/config"
	"user-service/internal/cache"
	"user-service/internal/dto"
	"user-service/internal/model"
	"user-service/internal/service"
)

type AuthHandler struct {
	userService service.UserService
	mailService service.MailService
}

func NewAuthHandler(u service.UserService, m service.MailService) *AuthHandler {
	return &AuthHandler{u, m}
}

var ctx = context.Background()

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"username": user.Name,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
	}

	// Cache the token in Redis
	err = cache.SetTokenInCache(signedToken, time.Hour*24)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not cache token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.userService.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create user"})
	}

	token, err := h.mailService.GenerateEmailOTP(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate email OTP"})
	}

	verifyLink := "http://your-frontend.com/verify?token=" + token
	go h.mailService.SendVerifyMail(user.Email, verifyLink)

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Get token from request
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing Authorization Header"})
	}

	tokenString := authHeader[len("Bearer "):]

	// Blacklist the token in Redis
	err := cache.BlacklistToken(tokenString, time.Hour*24)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not blacklist token"})
	}

	// Remove token from cache (optional, to ensure no further usage)
	err = cache.RedisClient.Del(ctx, tokenString).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not remove token from cache"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out successfully"})
}
func (h *AuthHandler) InActiveUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	err = h.userService.InActiveUser(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User InActive successfully"})
}
