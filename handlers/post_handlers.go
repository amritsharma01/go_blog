package handlers

import (
	"crud_api/errors"
	"crud_api/models"
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// func HandleAppError(c echo.Context, err error, fallbackMessage string) error {
// 	if appErr, ok := err.(*errors.AppError); ok {
// 		log.Printf("[ERROR] %s | Internal: %v", appErr.Message, appErr.Err)
// 		return utils.ErrorResponse(c, appErr.Code, appErr.Message)
// 	}
// 	return utils.ErrorResponse(c, http.StatusInternalServerError, fallbackMessage)
// }

type PostHandler struct {
	service services.PostService
}

func NewPostHandler(service services.PostService) *PostHandler {
	return &PostHandler{service}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new blog post (requires authentication)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post body requestmodels.CreatePostRequest true "Post content"
// @Success 201 {object} responsemodels.JSONResponseStruct{data=responsemodels.PostResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/posts [post]
func (h *PostHandler) CreatePost(c echo.Context) error {
	var req requestmodels.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(c,
			errors.BadRequest(
				"Invalid request body",
				"Failed to bind request body",
				err,
			),
			"",
		)
	}

	req.Sanitize()

	if req.Title == "" || req.Description == "" {
		return errors.HandleError(c,
			errors.BadRequest(
				"Category name is required",
				"Client sent empty category name",
				nil,
			),
			"",
		)
	}
	authUser := c.Get("user").(models.User)
	post := requestmodels.FromCreatePostRequest(req, authUser.ID)

	if err := h.service.Create(&post); err != nil {
		return errors.HandleError(c, err, "")
	}

	createdPost, err := h.service.GetByID(post.ID)
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	return responsemodels.JSONResponse(c, http.StatusCreated, "Successfully created post", responsemodels.ToPostResponse(*createdPost))
}

// GetPosts godoc
// @Summary Get list of posts
// @Description Get paginated list of posts with optional filters
// @Tags posts
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param category_id query string false "Filter by category ID"
// @Param author_id query string false "Filter by author ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} responsemodels.PaginatedResponse{data=[]responsemodels.PostResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/posts [get]
func (h *PostHandler) GetPosts(c echo.Context) error {
	search := c.QueryParam("search")
	categoryID := c.QueryParam("category_id")
	authorID := c.QueryParam("author_id")
	p := responsemodels.GetPagination(c)

	posts, total, err := h.service.GetAll(search, categoryID, authorID, p.Offset, p.Limit)
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	var response []responsemodels.PostResponse
	for _, post := range posts {
		response = append(response, responsemodels.ToPostResponse(post))
	}

	paginated := responsemodels.NewPaginatedResponse(response, p.Page, p.Limit, total)
	return responsemodels.SendPaginatedResponse(c, http.StatusOK, "Posts retrieved successfully", paginated)
}

// PostDetails godoc
// @Summary Get post details
// @Description Get details of a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} responsemodels.JSONResponseStruct{data=responsemodels.PostResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/posts/{id} [get]
func (h *PostHandler) PostDetails(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	return responsemodels.JSONResponse(c, http.StatusOK, "Post retrieved successfully", responsemodels.ToPostResponse(*post))
}

// PostDelete godoc
// @Summary Delete a post
// @Description Delete a post (only by author or admin)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} responsemodels.JSONResponseStruct
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/posts/{id} [delete]
func (h *PostHandler) PostDelete(c echo.Context) error {
	authUser := c.Get("user").(models.User)
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	if err := h.service.Delete(post, authUser.ID); err != nil {
		return errors.HandleError(c, err, "")
	}

	return responsemodels.JSONResponse(c, http.StatusOK, "Post deleted successfully", nil)
}

// PostEdit godoc
// @Summary Update a post
// @Description Update an existing post (only by author)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Param post body requestmodels.UpdatePostRequest true "Updated post content"
// @Success 200 {object} responsemodels.JSONResponseStruct{data=responsemodels.PostResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/posts/{id} [patch]
func (h *PostHandler) PostEdit(c echo.Context) error {
	authUser := c.Get("user").(models.User)
	id, _ := strconv.Atoi(c.Param("id"))

	var req requestmodels.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(c,
			errors.BadRequest(
				"Invalid request body",
				"Failed to bind request body",
				err,
			),
			"",
		)
	}

	if req.Title == "" || req.Description == "" {
		return errors.HandleError(c,
			errors.BadRequest(
				"Post title and description is required",
				"Client sent empty post title or description",
				nil,
			),
			"",
		)
	}

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	requestmodels.FromUpdatePostRequest(post, req)

	if err := h.service.Update(post, authUser.ID); err != nil {
		return errors.HandleError(c, err, "")
	}

	updatedPost, err := h.service.GetByID(post.ID)
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	return responsemodels.JSONResponse(c, http.StatusOK, "Post updated successfully", responsemodels.ToPostResponse(*updatedPost))
}

// GetPostsbyAuthor godoc
// @Summary Get posts by author
// @Description Get paginated list of posts by specific author
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param author_id path int true "Author ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} responsemodels.PaginatedResponse{data=[]responsemodels.PostResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/authors/{author_id}/posts [get]
func (h *PostHandler) GetPostsbyAuthor(c echo.Context) error {
	authorID := c.Param("author_id")
	p := responsemodels.GetPagination(c)

	posts, total, err := h.service.GetByAuthorID(authorID, p.Offset, p.Limit)
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	var response []responsemodels.PostResponse
	for _, post := range posts {
		response = append(response, responsemodels.ToPostResponse(post))
	}

	paginated := responsemodels.NewPaginatedResponse(response, p.Page, p.Limit, total)
	return responsemodels.SendPaginatedResponse(c, http.StatusOK, "Posts retrieved successfully", paginated)
}
