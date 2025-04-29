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
	jwtMiddleware := middleware.NewJWTMiddleware(db)
	protected := e.Group("")
	protected.Use(jwtMiddleware.Middleware)

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Blog API!")
	})

	// Auth routes
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	e.POST("/auth/register", userHandler.Register)
	e.POST("/auth/login", userHandler.Login)

	// User routes (protected)
	protected.GET("/users", userHandler.GetAllUsers)

	// Post routes
	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	e.GET("/posts", postHandler.GetPosts)        // Public paginated post listing
	e.GET("/posts/:id", postHandler.PostDetails) // Public post details by ID

	protected.POST("/posts", postHandler.CreatePost)                         // Create post
	protected.PATCH("/posts/:id", postHandler.PostEdit)                      // Update post
	protected.DELETE("/posts/:id", postHandler.PostDelete)                   // Delete post
	protected.GET("/authors/:author_id/posts", postHandler.GetPostsbyAuthor) // Posts by specific author

	// Category routes
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	protected.GET("/categories", categoryHandler.ListCategories)        // Paginated list
	protected.POST("/categories", categoryHandler.AddCategory)          // Create
	protected.DELETE("/categories/:id", categoryHandler.DeleteCategory) // Delete
}
