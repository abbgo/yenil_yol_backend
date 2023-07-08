package helpers

import (
	"regexp"
)

func ValidatePhoneNumber(regPersonalNumber string) bool {
	regexpPersonalNumber := regexp.MustCompile(`^(\+9936)[1-5][0-9]{6}$`)
	isMatchPersonalNumber := regexpPersonalNumber.MatchString(regPersonalNumber)
	return isMatchPersonalNumber
}
