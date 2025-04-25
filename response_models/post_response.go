package responsemodels

type PostResponse struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Author      AuthorInfo   `json:"author"`
	Category    CategoryInfo `json:"category"`
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
