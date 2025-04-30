// @title CRUD API
// @version 1.0
// @description This is a sample CRUD API with JWT authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"crud_api/config"
	"log"

	"crud_api/routes"

	_ "crud_api/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	e.Use(middleware.CORS())

	// Register routes
	routes.RegisterRoutes(e, db)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Start the server
	e.Logger.Fatal(e.Start(":8000"))
}
