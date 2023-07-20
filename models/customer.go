package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type Customer struct {
	ID          string `json:"id,omitempty"`
	FullName    string `json:"full_name,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password,omitempty" binding:"required"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

func ValidateCustomer(phoneNumber, customerID string, isRegisterFunction bool) error {

	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if isRegisterFunction {
		var phone_number string
		db.QueryRow(context.Background(), "SELECT phone_number FROM customers WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&phone_number)
		if phone_number != "" {
			return errors.New("this customer already exists")
		}
	} else {
		if customerID == "" {
			return errors.New("customer_id is required")
		}

		// database - de request body - den gelen id bilen gabat gelyan customer barmy ya-da yokmy sol barlanyar
		// eger yok bolsa onda error return edilyar
		var id string
		if err := db.QueryRow(context.Background(), "SELECT id FROM customers WHERE id = $1 AND deleted_at IS NULL", customerID).Scan(&id); err != nil {
			return errors.New("customer not found")
		}

		var customer_id string
		db.QueryRow(context.Background(), "SELECT id FROM customers WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&customer_id)
		if customer_id != customerID && customer_id != "" {
			return errors.New("this customer already exists")
		}

	}

	if !helpers.ValidatePhoneNumber(phoneNumber) {
		return errors.New("invalid phone number")
	}

	return nil

}
