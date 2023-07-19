package models

type PageTranslation struct {
	ID            string `json:"id,omitempty"`
	TitleTM       string `json:"title_tm,omitempty"`
	TitleRU       string `json:"title_ru,omitempty"`
	DescriptionTM string `json:"description_tm,omitempty"`
	DescriptionRU string `json:"description_ru,omitempty"`
	PageID        string `json:"page_id,omitempty" binding:"required"`
	CreatedAt     string `json:"-"`
	UpdatedAt     string `json:"-"`
	DeletedAt     string `json:"-"`
}

type PageTranslationUpdate struct {
	ID            string `json:"id,omitempty" binding:"required"`
	TitleTM       string `json:"title_tm,omitempty"`
	TitleRU       string `json:"title_ru,omitempty"`
	DescriptionTM string `json:"description_tm,omitempty"`
	DescriptionRU string `json:"description_ru,omitempty"`
	PageID        string `json:"page_id,omitempty" binding:"required"`
}
