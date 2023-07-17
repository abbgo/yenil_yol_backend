package models

type PageTranslation struct {
	ID            string `json:"id,omitempty"`
	TitleTM       string `json:"title_tm,omitempty" binding:"required"`
	TitleRU       string `json:"title_ru,omitempty" binding:"required"`
	DescriptionTM string `json:"description_tm,omitempty" binding:"required"`
	DescriptionRU string `json:"description_ru,omitempty" binding:"required"`
	PageID        string `json:"page_id,omitempty" binding:"required"`
	CreatedAt     string `json:"-"`
	UpdatedAt     string `json:"-"`
	DeletedAt     string `json:"-"`
}
