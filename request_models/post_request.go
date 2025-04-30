package requestmodels

// Mappers to models.Post
import (
	"crud_api/models"
	"strings"
)

type CreatePostRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	CategoryID  uint   `json:"category_id,omitempty"` // optional
}

type UpdatePostRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	CategoryID  uint   `json:"category_id,omitempty"` // optional
}

func FromCreatePostRequest(req CreatePostRequest, authorID uint) models.Post {
	return models.Post{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID,
	}
}

func FromUpdatePostRequest(post *models.Post, req UpdatePostRequest) {
	post.Title = req.Title
	post.Description = req.Description
	post.CategoryID = req.CategoryID
	post.Category = models.Category{}
}

func (r *CreatePostRequest) Sanitize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
}
