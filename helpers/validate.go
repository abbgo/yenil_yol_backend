package helpers

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
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

func ValidateShopOwnerByToken(c *gin.Context, db *pgxpool.Pool, shopOwnerID string) error {
	// middleware - den gelen id - ni alyar
	ID, hasID := c.Get("id")
	if !hasID {
		return errors.New("id is required")
	}
	var ok bool
	id, ok := ID.(string)
	if !ok {
		return errors.New("id must be string")
	}

	// bu yerde ilki bilen tokenden gelen id admina degisliligi barlanyar
	// ol tokenden gelen id admina degisli dal bolsa , onda onun shop_owner - e
	// degisliligi barlanyar , eger onada degisli dal bolsa error return edilyar
	if err := ValidateRecordByID("admins", id, "NULL", db); err != nil {
		if id != shopOwnerID {
			return errors.New("this shop isn't for you")
		}
	}

	return nil
}
