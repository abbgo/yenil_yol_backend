package models

type Dimension struct {
	ID        string `json:"id,omitempty"`
	Dimension string `json:"dimension" binding:"required"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}
