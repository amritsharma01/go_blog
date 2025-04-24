package handlers

import (
	"crud_api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

// Constructor
func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var post models.Post

	// Bind the request body to the Post struct
	if err := c.Bind(&post); err != nil {
		// Log the error message
		c.Logger().Errorf("Error binding request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	// Check if all required fields are filled
	if post.Title == "" || post.Description == "" || post.AuthorID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Title, Description, and AuthorID are required"})
	}

	// Check email uniqueness else error
	var existing models.Post
	if err := h.DB.Where("title = ? AND author_id = ?", post.Title, post.AuthorID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "Post already exists"})
	}

	// Check if the author exists in the database
	var author models.User
	if err := h.DB.First(&author, post.AuthorID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Author not found"})
	}

	// Save post
	if err := h.DB.Create(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create post"})
	}
	// Preload the Author data to include in the response
	if err := h.DB.Preload("Author").First(&post, post.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve author"})
	}

	// Preload the category data to include in the response
	if err := h.DB.Preload("Category").First(&post, post.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve category"})
	}

	return c.JSON(http.StatusCreated, post)
}

// List all posts with pagination
func (h *PostHandler) GetPosts(c echo.Context) error {
	var posts []models.Post

	// Get pagination parameters from the query string
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	// Default values for pagination
	pageNum := 1
	limitNum := 10

	// Parse page and limit if provided
	if page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err == nil && parsedPage > 0 {
			pageNum = parsedPage
		}
	}

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 {
			limitNum = parsedLimit
		}
	}

	// Calculate the offset based on the current page
	offset := (pageNum - 1) * limitNum

	// Preload the Author data for each post
	if err := h.DB.Limit(limitNum).Offset(offset).Preload("Author").Find(&posts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve posts"})
	}

	// Optionally, return total count for pagination metadata
	var totalPosts int64
	if err := h.DB.Model(&models.Post{}).Count(&totalPosts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to count posts"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       posts,
		"page":       pageNum,
		"limit":      limitNum,
		"total":      totalPosts,
		"totalPages": int(totalPosts / int64(limitNum)),
	})
}

func (h *PostHandler) PostDetails(c echo.Context) error {
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	var post models.Post

	// Find the post by ID and preload the Author data
	if err := h.DB.Preload("Author").First(&post, postID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Post not found"})
	}

	// Preload the category data to include in the response
	if err := h.DB.Preload("Category").First(&post, post.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve category"})
	}

	return c.JSON(http.StatusOK, post)
}

func (h *PostHandler) PostDelete(c echo.Context) error {
	// Get the post ID from the URL parameter
	postID := c.Param("id")

	// Delete the post with the given ID
	if err := h.DB.Delete(&models.Post{}, postID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete post"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Post deleted successfully"})
}

func (h *PostHandler) PostEdit(c echo.Context) error {
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	var post models.Post

	// Bind the request body to the post struct
	if err := c.Bind(&post); err != nil {
		c.Logger().Errorf("Error binding request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	// Find the post by ID
	if err := h.DB.First(&post, postID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Post not found"})
	}

	// Update the post's title and description
	post.Title = post.Title
	post.Description = post.Description

	// Save the updated post
	if err := h.DB.Save(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update post"})
	}

	return c.JSON(http.StatusOK, post)
}

// Get posts by author with pagination
func (h *PostHandler) GetPostsbyAuthor(c echo.Context) error {
	authorID := c.Param("author_id")

	// Convert author_id to uint (assuming it's a numeric value)
	authorIDUint, err := strconv.Atoi(authorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid author ID",
		})
	}

	// Get pagination parameters from the query string
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	// Default values for pagination
	pageNum := 1
	limitNum := 10

	// Parse page and limit if provided
	if page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err == nil && parsedPage > 0 {
			pageNum = parsedPage
		}
	}

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 {
			limitNum = parsedLimit
		}
	}

	// Calculate the offset based on the current page
	offset := (pageNum - 1) * limitNum

	// Retrieve the posts by the specific author
	var posts []models.Post
	if err := h.DB.Limit(limitNum).Offset(offset).Preload("Author").Where("author_id = ?", authorIDUint).Find(&posts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve posts",
		})
	}

	// Optionally, return total count for pagination metadata
	var totalPosts int64
	if err := h.DB.Model(&models.Post{}).Where("author_id = ?", authorIDUint).Count(&totalPosts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to count posts"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       posts,
		"page":       pageNum,
		"limit":      limitNum,
		"total":      totalPosts,
		"totalPages": int(totalPosts / int64(limitNum)),
	})
}
