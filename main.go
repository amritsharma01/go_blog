package main

import (
	"crud_api/config"

	"crud_api/routes"
	"crud_api/validators"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Connect to the database
	db := config.ConnectDB()

	// Initialize validators
	validators.InitValidators(db)

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
