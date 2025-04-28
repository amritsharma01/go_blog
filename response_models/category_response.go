package responsemodels

import "crud_api/models"

type CategoryResponse struct {
	ID   uint   `json:"cid"`
	Name string `json:"cname"`
}

func ToCatResponse(c models.Category) CategoryResponse {
	return CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}
