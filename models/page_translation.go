package models

type PageTranslation struct {
	ID            string `json:"id,omitempty"`
	TextTitleTM   string `json:"text_title_tm,omitempty"`
	TextTitleRU   string `json:"text_title_ru,omitempty"`
	DescriptionTM string `json:"description_tm,omitempty"`
	DescriptionRU string `json:"description_ru,omitempty"`
	PageID        string `json:"page_id,omitempty" binding:"required"`
	OrderNumber   uint   `json:"order_number,omitempty"`
	CreatedAt     string `json:"-"`
	UpdatedAt     string `json:"-"`
	DeletedAt     string `json:"-"`
}

type PageTranslationUpdate struct {
	ID            string `json:"id,omitempty" binding:"required"`
	TextTitleTM   string `json:"text_title_tm,omitempty"`
	TextTitleRU   string `json:"text_title_ru,omitempty"`
	DescriptionTM string `json:"description_tm,omitempty"`
	DescriptionRU string `json:"description_ru,omitempty"`
	PageID        string `json:"page_id,omitempty" binding:"required"`
	OrderNumber   uint   `json:"order_number,omitempty"`
}
