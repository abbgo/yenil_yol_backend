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
	Password    string `json:"password,omitempty"`
}

type CustomerUpdatePassword struct {
	ID       string `json:"id,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

func ValidateCustomer(customer Customer, isRegisterFunction bool) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if isRegisterFunction {
		if customer.Password == "" {
			return errors.New("password is required")
		}

		var phone_number string
		db.QueryRow(context.Background(), "SELECT phone_number FROM customers WHERE phone_number = $1 AND deleted_at IS NULL", customer.PhoneNumber).Scan(&phone_number)
		if phone_number != "" {
			return errors.New("this customer already exists")
		}
	} else {
		if customer.ID == "" {
			return errors.New("customer_id is required")
		}

		if err := helpers.ValidateRecordByID("customers", customer.ID, "NULL", db); err != nil {
			return err
		}

		var customer_id string
		db.QueryRow(context.Background(), "SELECT id FROM customers WHERE phone_number = $1 AND deleted_at IS NULL", customer.PhoneNumber).Scan(&customer_id)
		if customer_id != customer.ID && customer_id != "" {
			return errors.New("this customer already exists")
		}

	}

	if !helpers.ValidatePhoneNumber(customer.PhoneNumber) {
		return errors.New("invalid phone number")
	}

	return nil
}
