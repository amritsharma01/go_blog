package responsemodels

import (
	"crud_api/models"
	"time"
)

type PostResponse struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Author      AuthorInfo   `json:"author"`
	Category    CategoryInfo `json:"category"`
	Created     string       `json:"created_at"`
}

type AuthorInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToPostResponse(p models.Post) PostResponse {
	return PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Created:     p.CreatedAt.Format(time.RFC3339),

		Author: AuthorInfo{
			ID:    p.Author.ID,
			Name:  p.Author.Name,
			Email: p.Author.Email,
		},
		Category: CategoryInfo{
			ID:   p.Category.ID,
			Name: p.Category.Name,
		},
	}
}
