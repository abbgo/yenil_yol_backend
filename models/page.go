package models

type Page struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty" binding:"required"`
	Image     string `json:"image,omitempty"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type PageUpdate struct {
	ID    string `json:"id,omitempty" binding:"required"`
	Name  string `json:"name,omitempty" binding:"required"`
	Image string `json:"image,omitempty"`
}
