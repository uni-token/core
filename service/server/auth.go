package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"uni-token-service/logic"
	"uni-token-service/store"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Status   string `json:"status"` // "success", "error", "not_registered"
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}

func SetupAuthAPI(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", handleUserLogin)
		auth.POST("/register", handleRegister)
	}
}

func handleUserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Status:  "error",
			Message: "Invalid request data",
		})
		return
	}

	users, err := store.Users.List()
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusOK, AuthResponse{
			Status: "not_registered",
		})
		return
	}

	user, err := store.Users.Get(req.Username)
	if err != nil {
		c.JSON(http.StatusOK, AuthResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	hashedPassword := logic.HashPassword(req.Password)
	if user.Password != hashedPassword {
		c.JSON(http.StatusOK, AuthResponse{
			Status:  "error",
			Message: "Invalid username or password",
		})
		return
	}

	token, err := logic.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Status:  "error",
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Status:   "success",
		Username: user.Username,
		Token:    token,
	})
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Status:  "error",
			Message: "Invalid request data",
		})
		return
	}

	_, err := store.Users.Get(req.Username)
	if err == nil {
		c.JSON(http.StatusOK, AuthResponse{
			Status:  "error",
			Message: "User already exists",
		})
		return
	}

	hashedPassword := logic.HashPassword(req.Password)
	user := store.UserInfo{
		Username: req.Username,
		Password: hashedPassword,
	}

	err = store.Users.Put(req.Username, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Status:  "error",
			Message: "Failed to register",
		})
		return
	}

	token, err := logic.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Status:  "error",
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Status:   "success",
		Username: user.Username,
		Token:    token,
	})
}

func RequireUserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := logic.ValidateJWT(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Id)
		c.Next()
	}
}
