package helpers

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ValidatePhoneNumber(regPersonalNumber string) bool {
	regexpPersonalNumber := regexp.MustCompile(`^(\+9936)[1-5][0-9]{6}$`)
	isMatchPersonalNumber := regexpPersonalNumber.MatchString(regPersonalNumber)
	return isMatchPersonalNumber
}

func ValidateEmailAddress(email string) bool {
	// Regular expression pattern for email addresses
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the email matches the pattern
	return regex.MatchString(email)
}

func ValidateRecordByID(tableName, id, nullStr string, db *pgxpool.Pool) error {
	// database - de request gelen id bilen gabat gelyan maglumat barmy ya-da yokmy sol barlanyar
	// eger yok bolsa onda error return edilyar
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = '%s' AND deleted_at IS %s", tableName, id, nullStr)
	fmt.Println(query)
	if err := db.QueryRow(context.Background(), query).Scan(&id); err != nil {
		return errors.New("record not found")
	}
	return nil
}

func ValidateStructData(s interface{}) error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}
	return nil
}
