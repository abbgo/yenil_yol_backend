package models

type Brend struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name" binding:"required"`
	Image string `json:"image,omitempty" binding:"required"`
	Slug  string `json:"slug,omitempty"`
}
