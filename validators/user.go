package validators

import (
	"github.com/anthonyhawkins/savorbook/forms"
	"github.com/go-playground/validator/v10"
)

func ValidateRegistration(registration forms.RegisterForm) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(registration)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.value = err.Param()
			errors = append(errors, &element)
		}
	}

	if registration.Password != registration.PasswordConfirm {
		var element ErrorResponse
		element.FailedField = "RegistrationForm.PasswordConfirm"
		element.Message = "Password and Confirm Password must match."
		errors = append(errors, &element)
	}

	return errors
}
