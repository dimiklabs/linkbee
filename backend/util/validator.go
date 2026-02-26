package util

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

const (
	MinPasswordLength     = 8
	MaxPasswordLength     = 128
	MinComplexityCriteria = 3
)

func ValidatePasswordComplexity(password string) bool {
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	criteriaCount := 0
	if hasUpper {
		criteriaCount++
	}
	if hasLower {
		criteriaCount++
	}
	if hasNumber {
		criteriaCount++
	}
	if hasSpecial {
		criteriaCount++
	}

	return criteriaCount >= MinComplexityCriteria
}

func PasswordComplexityValidator(fl validator.FieldLevel) bool {
	return ValidatePasswordComplexity(fl.Field().String())
}
