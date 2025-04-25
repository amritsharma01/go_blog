package handlers

import (
	"crud_api/models"
	responsemodels "crud_api/response_models"
	"crud_api/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ToPostResponse(p models.Post) responsemodels.PostResponse {
	return responsemodels.PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Author: responsemodels.AuthorInfo{
			ID:    p.Author.ID,
			Name:  p.Author.Name,
			Email: p.Author.Email,
		},
		Category: responsemodels.CategoryInfo{
			ID:   p.Category.ID,
			Name: p.Category.Name,
		},
	}
}

type PostHandler struct {
	DB *gorm.DB
}

// Constructor
func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req models.Post

	// Bind DTO
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body")
	}

	// Basic validation
	if req.Title == "" || req.Description == "" || req.AuthorID == 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Title, Description and AuthorID are required")
	}

	// Check post uniqueness
	var existing models.Post
	if err := h.DB.Where("title = ? AND author_id = ?", req.Title, req.AuthorID).First(&existing).Error; err == nil {
		return utils.ErrorResponse(c, http.StatusConflict, "Post already exists")
	}

	// Check author existence
	var author models.User
	if err := h.DB.First(&author, req.AuthorID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Author not found")
	}

	// Optional: validate category
	if req.CategoryID != nil {
		var cat models.Category
		if err := h.DB.First(&cat, *req.CategoryID).Error; err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Category not found")
		}
	}

	// Save
	if err := h.DB.Create(&req).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
	}

	// Preload for response
	if err := h.DB.Preload("Author").Preload("Category").First(&req, req.ID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load full post data")
	}

	// Return DTO
	return utils.JSONResponse(c, http.StatusCreated, "Successfully created post", ToPostResponse(req))
}

func (h *PostHandler) GetPosts(c echo.Context) error {
	var posts []models.Post

	// --- Pagination ---
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pageNum := 1
	limitNum := 10

	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitNum = l
	}

	offset := (pageNum - 1) * limitNum

	// --- Fetch posts with Author and Category ---
	if err := h.DB.Limit(limitNum).Offset(offset).Preload("Author").Preload("Category").Find(&posts).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve posts")
	}

	// --- Count total posts for pagination info ---
	var totalPosts int64
	if err := h.DB.Model(&models.Post{}).Count(&totalPosts).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to count posts")
	}

	// --- Convert to DTO (PostResponse) ---
	var responsePosts []responsemodels.PostResponse
	for _, post := range posts {
		responsePosts = append(responsePosts, ToPostResponse(post))
	}

	paginatedResponse := echo.Map{
		"posts":      responsePosts,
		"page":       pageNum,
		"limit":      limitNum,
		"total":      totalPosts,
		"totalPages": int((totalPosts + int64(limitNum) - 1) / int64(limitNum)), // rounding up
	}

	// --- Return clean response ---
	return utils.JSONResponse(c, http.StatusOK, "Posts retrieved successfully", paginatedResponse)
}

// get the post details
func (h *PostHandler) PostDetails(c echo.Context) error {
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	var post models.Post

	// Find the post by ID and preload the Author data
	if err := h.DB.Preload("Author").First(&post, postID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")

	}

	// Preload the category data to include in the response
	if err := h.DB.Preload("Category").First(&post, post.ID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve category")

	}
	return utils.JSONResponse(c, http.StatusOK, "Post Retrieved Succesfully", ToPostResponse(post))

}

// delete the posts
func (h *PostHandler) PostDelete(c echo.Context) error {
	// Get the post ID from the URL parameter
	postID := c.Param("id")

	// Check if the post exists before deleting
	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	// Delete the post with the given ID
	if err := h.DB.Delete(&post).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post")
	}

	return utils.JSONResponse(c, http.StatusOK, "Post deleted successfully", nil)
}

// edit the post
func (h *PostHandler) PostEdit(c echo.Context) error {
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	var req models.Post

	// Bind the request body to the post struct
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Title == "" || req.Description == "" || req.AuthorID == 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Title, Description and AuthorID are required")
	}

	// Find the post by ID
	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	// Update the post with the new values
	post.Title = req.Title
	post.Description = req.Description
	post.AuthorID = req.AuthorID
	post.CategoryID = req.CategoryID

	// Save the updated post
	if err := h.DB.Save(&post).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post")
	}

	// Preload the Author and Category for the response
	if err := h.DB.Preload("Author").Preload("Category").First(&post, post.ID).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load full post data")
	}

	// Return DTO
	return utils.JSONResponse(c, http.StatusOK, "Successfully updated post", ToPostResponse(post))
}

// get posts according to author
func (h *PostHandler) GetPostsbyAuthor(c echo.Context) error {
	authorID := c.Param("author_id")

	// Convert author_id to uint (assuming it's a numeric value)
	authorIDUint, err := strconv.Atoi(authorID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid author ID")
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
	if err := h.DB.Limit(limitNum).Offset(offset).Preload("Author").Preload("Category").Where("author_id = ?", authorIDUint).Find(&posts).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve posts")
	}

	// Optionally, return total count for pagination metadata
	var totalPosts int64
	if err := h.DB.Model(&models.Post{}).Where("author_id = ?", authorIDUint).Count(&totalPosts).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to count posts")
	}

	// Convert posts to DTOs
	var responsePosts []responsemodels.PostResponse
	for _, post := range posts {
		responsePosts = append(responsePosts, ToPostResponse(post))
	}

	// Paginated response
	paginatedResponse := echo.Map{
		"posts":      responsePosts,
		"page":       pageNum,
		"limit":      limitNum,
		"total":      totalPosts,
		"totalPages": int((totalPosts + int64(limitNum) - 1) / int64(limitNum)),
	}

	// Return the paginated response
	return utils.JSONResponse(c, http.StatusOK, "Posts retrieved successfully", paginatedResponse)
}
