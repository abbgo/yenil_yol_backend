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
}

func ValidateSetting(setting Setting) error {
	if !helpers.ValidatePhoneNumber(setting.PhoneNumber) {
		return errors.New("invalid phone number")
	}

	if !helpers.ValidateEmailAddress(setting.Email) {
		return errors.New("invalid mail address")
	}

	return nil
}
