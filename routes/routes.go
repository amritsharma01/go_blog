package routes

import (
	"crud_api/handlers"
	"crud_api/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)
	categoryHandler := handlers.NewCategory(db)
	jwtMiddleware := middleware.NewJWTMiddleware(db)

	// Public routes (no authentication needed)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Blog API!")
	})
	e.POST("/auth/register", userHandler.CreateUser)
	e.POST("/auth/login", userHandler.Login)

	// Protected routes (require JWT)
	protected := e.Group("")
	protected.Use(jwtMiddleware.Middleware)
	{
		protected.POST("/posts/create", postHandler.CreatePost)
		protected.GET("/users", userHandler.GetUsers)
		protected.GET("/posts/list", postHandler.GetPosts)
		protected.GET("/posts/list/:id", postHandler.PostDetails)
		protected.DELETE("/posts/delete/:id", postHandler.PostDelete)
		protected.PATCH("/posts/edit/:id", postHandler.PostEdit)
		protected.GET("/posts/author/:author_id", postHandler.GetPostsbyAuthor)
		protected.POST("/category/add", categoryHandler.AddCategory)
		protected.GET("/category/list", categoryHandler.ListCategories)
	}
}
