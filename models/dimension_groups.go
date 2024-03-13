package models

type DimensionGroup struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name" binding:"required"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}
