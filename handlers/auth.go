package handlers

import (
	"api/database"
	"api/models"
	"net/http"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Name == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and password are required"})
		return
	}

	var user models.User
	if err := database.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

func Register(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = user.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	database.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
