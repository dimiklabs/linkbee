package util

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/shafikshaon/linkbee/constant"
)

func TranslateValidationError(err error) (string, string) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := strings.ToLower(fieldErr.Field())
			tag := fieldErr.Tag()

			switch tag {
			case "password_complexity":
				return constant.ErrCodeWeakPassword, constant.ErrMsgWeakPassword
			case "required":
				return constant.ErrCodeValidationError, field + " is required"
			case "email":
				return constant.ErrCodeValidationError, "Invalid email format"
			case "min":
				return constant.ErrCodeValidationError, field + " must be at least " + fieldErr.Param() + " characters"
			case "max":
				return constant.ErrCodeValidationError, field + " must be at most " + fieldErr.Param() + " characters"
			default:
				return constant.ErrCodeValidationError, "Invalid " + field
			}
		}
	}
	return constant.ErrCodeValidationError, err.Error()
}
