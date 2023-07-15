package models

type Setting struct {
	ID          string `json:"id,omitempty"`
	Logo        string `json:"logo,omitempty" binding:"required"`
	Favicon     string `json:"favicon,omitempty" binding:"required"`
	Email       string `json:"email,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}
