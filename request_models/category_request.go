package requestmodels

import "crud_api/models"

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

func FromCatRequest(req CategoryRequest) models.Category {
	return models.Category{
		Name: req.Name,
	}
}
