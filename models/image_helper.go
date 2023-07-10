package models

type ImageHelper struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty" binding:"required"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}
