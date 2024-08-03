package models

type Complaint struct {
	ID     string `json:"id,omitempty"`
	TextTM string `json:"txt_tm,omitempty" binding:"required"`
	TextRU string `json:"txt_ru,omitempty" binding:"required"`
}
