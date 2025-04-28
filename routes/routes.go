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

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {

	//using jwt middleware``
	jwtMiddleware := middleware.NewJWTMiddleware(db)

	// Protected routes (require JWT)
	protected := e.Group("")
	protected.Use(jwtMiddleware.Middleware)

	//common public route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Blog API!")
	})

	//dependency injection for user
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	//public routes for users
	e.POST("/auth/register", userHandler.Register)
	e.POST("/auth/login", userHandler.Login)

	//private routes for users
	protected.GET("/users", userHandler.GetAllUsers)

	//dependency injection for posts
	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	//public routes for posts
	e.GET("/posts/list", postHandler.GetPosts)
	e.GET("/posts/list/:id", postHandler.PostDetails)

	//protected routes for posts
	protected.POST("/posts/create", postHandler.CreatePost)
	protected.DELETE("/posts/delete/:id", postHandler.PostDelete)
	protected.PATCH("/posts/edit/:id", postHandler.PostEdit)
	protected.GET("/posts/author/:author_id", postHandler.GetPostsbyAuthor)

	// dependency injection for category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	//public routes for category
	protected.GET("/category/list", categoryHandler.ListCategories)

	//private routes for category
	protected.POST("/category/add", categoryHandler.AddCategory)
	protected.DELETE("/category/delete/:id", categoryHandler.DeleteCategory)

}
