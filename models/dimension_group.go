package models

type DimensionGroup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name" binding:"required"`
}
