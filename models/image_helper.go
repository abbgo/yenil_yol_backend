package models

type HelperImage struct {
	ID    string `json:"id,omitempty"`
	Image string `json:"image,omitempty" binding:"required"`
}
