package models

type Category struct {
	ID        string `json:"id,omitempty"`
	NameTM    string `json:"name_tm,omitempty" binding:"required"`
	NameRU    string `json:"name_ru,omitempty" binding:"required"`
	Image     string `json:"image,omitempty"`
	SlugTM    string `json:"slug_tm,omitempty"`
	SlugRU    string `json:"slug_ru,omitempty"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type CategoryUpdate struct {
	ID     string `json:"id,omitempty" binding:"required"`
	NameTM string `json:"name_tm,omitempty" binding:"required"`
	NameRU string `json:"name_ru,omitempty" binding:"required"`
	Image  string `json:"image,omitempty"`
}
