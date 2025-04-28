package requestmodels

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}
