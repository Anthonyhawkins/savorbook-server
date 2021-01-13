package users

import (
	"github.com/go-playground/validator/v10"
)

func ValidateRegistration(registration RegisterForm) []string {
	var errors []string
	validate := validator.New()
	err := validate.Struct(registration)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}

	return errors
}
