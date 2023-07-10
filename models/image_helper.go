package models

type HelperImage struct {
	ID        string `json:"id,omitempty"`
	Image     string `json:"image,omitempty" binding:"required"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}
