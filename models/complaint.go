package models

type Complaint struct {
	ID     string `json:"id,omitempty"`
	TextTM string `json:"text_tm,omitempty" binding:"required"`
	TextRU string `json:"text_ru,omitempty" binding:"required"`
}
