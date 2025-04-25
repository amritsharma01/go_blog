package main

import (
	"crud_api/config"
	"log"

	"crud_api/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// Connect to the database
	db := config.ConnectDB()

	// Create an echo instance
	e := echo.New()
	// Middleware for logging and recovery
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	routes.RegisterRoutes(e, db)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
