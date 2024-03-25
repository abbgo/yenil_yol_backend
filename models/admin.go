package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type Admin struct {
	ID           string `json:"id,omitempty"`
	FullName     string `json:"full_name,omitempty" binding:"required"`
	PhoneNumber  string `json:"phone_number,omitempty" binding:"required"`
	Password     string `json:"password,omitempty" binding:"required"`
	IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
}

type AdminUpdate struct {
	ID           string `json:"id,omitempty" binding:"required"`
	FullName     string `json:"full_name,omitempty" binding:"required"`
	PhoneNumber  string `json:"phone_number,omitempty" binding:"required"`
	IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
}

type AdminUpdatePassword struct {
	ID       string `json:"id,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

func ValidateAdmin(phoneNumber, adminID string, isRegisterFunction bool) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if isRegisterFunction {
		var phone_number string
		db.QueryRow(context.Background(), "SELECT phone_number FROM admins WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&phone_number)
		if phone_number != "" {
			return errors.New("this admin already exists")
		}
	} else {
		if adminID == "" {
			return errors.New("admin_id is required")
		}

		if err := helpers.ValidateRecordByID("admins", adminID, "NULL", db); err != nil {
			return err
		}

		var admin_id string
		db.QueryRow(context.Background(), "SELECT id FROM admins WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&admin_id)
		if admin_id != adminID && admin_id != "" {
			return errors.New("this admin already exists")
		}

	}

	if !helpers.ValidatePhoneNumber(phoneNumber) {
		return errors.New("invalid phone number")
	}

	return nil
}
