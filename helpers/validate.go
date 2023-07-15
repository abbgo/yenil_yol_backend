package helpers

import (
	"regexp"
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
