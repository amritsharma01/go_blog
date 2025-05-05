package routes

import (
	"crud_api/handlers"
	"crud_api/middleware"
	"crud_api/repositories"
	"crud_api/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterRoutes sets up all the routes for the application
func RegisterRoutes(e *echo.Echo, db *gorm.DB) {

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Blog API!")
	})

	// Auth routes
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	jwtMiddleware := middleware.NewJWTMiddleware(userService)
	protected := e.Group("")
	protected.Use(jwtMiddleware.Middleware)

	e.POST("/v1/auth/register", userHandler.Register)
	e.POST("/v1/auth/login", userHandler.Login)

	// User routes (protected)
	protected.GET("/v1/users", userHandler.GetAllUsers)

	// Post routes
	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	e.GET("/v1/posts", postHandler.GetPosts)        // Public paginated post listingo
	e.GET("/v1/posts/:id", postHandler.PostDetails) // Public post details by ID

	protected.POST("/v1/posts", postHandler.CreatePost)                         // Create post
	protected.PATCH("/v1/posts/:id", postHandler.PostEdit)                      // Update post
	protected.DELETE("/v1/posts/:id", postHandler.PostDelete)                   // Delete post
	protected.GET("/v1/authors/:author_id/posts", postHandler.GetPostsbyAuthor) // Posts by specific author

	// Category routes
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	protected.GET("/v1/categories", categoryHandler.ListCategories)        // Paginated list
	protected.POST("/v1/categories", categoryHandler.AddCategory)          // Create
	protected.DELETE("/v1/categories/:id", categoryHandler.DeleteCategory) // Delete
}
