package routes

import (
	"crud_api/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)
	categoryHandler := handlers.NewCategory(db)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Blog API!")
	})

	e.POST("/users", userHandler.CreateUser)
	e.POST("/posts/create", postHandler.CreatePost)
	e.GET("/users", userHandler.GetUsers)
	e.GET("/posts/list", postHandler.GetPosts)
	e.GET("/posts/list/:id", postHandler.PostDetails)
	e.DELETE("/posts/delete/:id", postHandler.PostDelete)
	e.PATCH("/posts/edit/:id", postHandler.PostEdit)
	e.GET("/posts/author/:author_id", postHandler.GetPostsbyAuthor)
	e.POST("/category/add", categoryHandler.AddCategory)
	//e.GET("category/list")

}
