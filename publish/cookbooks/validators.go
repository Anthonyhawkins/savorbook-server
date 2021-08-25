package cookbooks

import "github.com/go-playground/validator/v10"

type CookbookValidator struct {
	Cookbook struct {
		Title    string             `json:"title" validate:"max=75"`
		SubTitle string             `json:"subTitle" validate:"max=75"`
		Blurb    string             `json:"blurb" validate:"max=500"`
		Image    string             `json:"image" validate:"omitempty"`
		Sections []SectionValidator `json:"sections" validate:"dive"`
	} `json:"cookbook"`
	Model CookbookModel `json:"-"`
}

type SectionValidator struct {
	Name     string `json:"name" validate:"max=75"`
	Overview string `json:"overview" validate:"max=500"`
	Recipes  []uint `json:"recipes" validate:"dive,numeric"`
}

func NewCookbookValidator() *CookbookValidator {
	return &CookbookValidator{}
}

func (v *CookbookValidator) Validate() ([]string, error) {
	var errors []string
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " = " + err.Tag()
			errors = append(errors, message)
		}
	}

	return errors, err
}

func (v *CookbookValidator) BindModel(userID uint) error {
	v.Model.UserID = userID
	v.Model.Title = v.Cookbook.Title
	v.Model.SubTitle = v.Cookbook.SubTitle
	v.Model.Blurb = v.Cookbook.Blurb
	v.Model.Image = v.Cookbook.Image
	if err := v.Model.setSections(userID, v.Cookbook.Sections); err != nil {
		return err
	}
	return nil
}
