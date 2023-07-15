package models

type Category struct {
	ID        string `json:"id,omitempty"`
	NameTM    string `json:"name_tm,omitempty" binding:"required"`
	NameRU    string `json:"name_ru,omitempty" binding:"required"`
	Image     string `json:"image,omitempty"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}
