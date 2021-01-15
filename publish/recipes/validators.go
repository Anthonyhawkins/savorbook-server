package recipes

import (
	"github.com/go-playground/validator/v10"
)

func ValidateRecipe(recipe Recipe) []string {
	var errors []string
	validate := validator.New()

	err := validate.Struct(recipe)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}

	return errors
}
