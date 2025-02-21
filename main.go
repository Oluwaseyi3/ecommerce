package main

import (
	"api/database"
	"api/handlers"

	//	"ecommerce-api/handlers"
	//	"ecommerce-api/middleware"
	"api/models"

	"github.com/gin-gonic/gin"
	//
	// )
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	router := gin.Default()

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	router.Run(":8080")
}
