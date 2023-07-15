package models

import (
	"errors"
	"github/abbgo/yenil_yol/backend/helpers"
)

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

func ValidateSetting(phoneNumber, email string) error {

	if !helpers.ValidatePhoneNumber(phoneNumber) {
		return errors.New("invalid phone number")
	}

	if !helpers.ValidateEmailAddress(email) {
		return errors.New("invalid mail address")
	}

	return nil
}
