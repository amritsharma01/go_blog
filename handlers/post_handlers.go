package handlers

import (
	"crud_api/models"
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"
	"crud_api/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	service services.PostService
}

func NewPostHandler(service services.PostService) *PostHandler {
	return &PostHandler{service}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req requestmodels.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	authUser := c.Get("user").(models.User)
	post := requestmodels.FromCreatePostRequest(req, authUser.ID)

	if err := h.service.Create(&post); err != nil {
		if err == services.ErrPostAlreadyExists {
			return utils.ErrorResponse(c, http.StatusConflict, "Post already exists")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
	}

	createdPost, _ := h.service.GetByID(post.ID)

	return utils.JSONResponse(c, http.StatusCreated, "Successfully created post", responsemodels.ToPostResponse(*createdPost))
}

func (h *PostHandler) GetPosts(c echo.Context) error {
	search := c.QueryParam("search")
	categoryID := c.QueryParam("category_id")
	authorID := c.QueryParam("author_id")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	posts, total, err := h.service.GetAll(search, categoryID, authorID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve posts")
	}

	var response []responsemodels.PostResponse
	for _, p := range posts {
		response = append(response, responsemodels.ToPostResponse(p))
	}

	return utils.PaginatedResponse(c, http.StatusOK, "Posts retrieved successfully", response, page, limit, total)
}

func (h *PostHandler) PostDetails(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	return utils.JSONResponse(c, http.StatusOK, "Post retrieved successfully", responsemodels.ToPostResponse(*post))
}

func (h *PostHandler) PostDelete(c echo.Context) error {
	authUser := c.Get("user").(models.User)
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	if err := h.service.Delete(post, authUser.ID); err != nil {
		return utils.ErrorResponse(c, http.StatusForbidden, "You can only delete your own posts")
	}

	return utils.JSONResponse(c, http.StatusOK, "Post deleted successfully", nil)
}

func (h *PostHandler) PostEdit(c echo.Context) error {
	authUser := c.Get("user").(models.User)
	id, _ := strconv.Atoi(c.Param("id"))

	var req requestmodels.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	if post.AuthorID != authUser.ID {
		return utils.ErrorResponse(c, http.StatusForbidden, "You can only edit your own posts")
	}

	requestmodels.FromUpdatePostRequest(post, req)

	if err := h.service.Update(post); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post")
	}

	updatedPost, err := h.service.GetByID(post.ID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch updated post")
	}

	return utils.JSONResponse(c, http.StatusOK, "Post updated successfully", responsemodels.ToPostResponse(*updatedPost))
}

func (h *PostHandler) GetPostsbyAuthor(c echo.Context) error {
	authorID, _ := strconv.Atoi(c.Param("author_id"))

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	posts, total, err := h.service.GetByAuthorID(uint(authorID), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve posts")
	}

	var response []responsemodels.PostResponse
	for _, p := range posts {
		response = append(response, responsemodels.ToPostResponse(p))
	}

	return utils.PaginatedResponse(c, http.StatusOK, "Posts retrieved successfully", response, page, limit, total)
}
