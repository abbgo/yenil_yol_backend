package models

type Brend struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty" binding:"required"`
	Image string `json:"image,omitempty"`
	Slug  string `json:"slug,omitempty"`
}

type BrendUpdate struct {
	ID    string `json:"id,omitempty" binding:"required"`
	Name  string `json:"name,omitempty" binding:"required"`
	Image string `json:"image,omitempty"`
}
