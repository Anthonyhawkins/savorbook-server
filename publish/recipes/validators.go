package recipes

import "github.com/go-playground/validator/v10"

type RecipeValidator struct {
	Recipe struct {
		Name             string                      `json:"name"                validate:"required,max=75"`
		Image            string                      `json:"image"               validate:"omitempty"`
		Description      string                      `json:"description"         validate:"max=2600"`
		PrepTime         string                      `json:"prepTime"            validate:"max=120"`
		Servings         string                      `json:"servings"            validate:"max=120"`
		Tags             []string                    `json:"tags"                validate:"dive,alphanum"`
		DependentRecipes []RecipeDependencyValidator `json:"dependentRecipes"`
		IngredientGroups []IngredientGroupValidator  `json:"ingredientGroups"    validate:"required,dive"`
		Steps            []StepValidator             `json:"steps"               validate:"required,dive"`
	} `json:"recipe"`
	Model RecipeModel `json:"-"`
}

type RecipeDependencyValidator struct {
	DependentRecipe uint   `json:"id"`
	Qty             string `json:"qty"`
}

type IngredientGroupValidator struct {
	GroupName   string                `json:"groupName"   validate:"max=50"`
	Ingredients []IngredientValidator `json:"ingredients" validate:"required,dive"`
}

type IngredientValidator struct {
	Name string `json:"name" validate:"required,max=32"`
	Qty  string `json:"qty"  validate:"omitempty,max=6"`
	Unit string `json:"unit" validate:"omitempty,max=12"`
}

type StepValidator struct {
	Type       string               `json:"type"        validate:"oneof=text tipText imageLeft imageRight imageDouble imageTriple"`
	Text       string               `json:"text"        validate:"max=865"`
	StepImages []StepImageValidator `json:"images"      validate:"dive"`
}

type StepImageValidator struct {
	Image string `json:"src"`
	Text  string `json:"text"`
}

func NewRecipeValidator() *RecipeValidator {
	return &RecipeValidator{}
}

func (v *RecipeValidator) Validate() ([]string, error) {
	var errors []string
	validate := validator.New()

	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}

	return errors, err
}

func (v *RecipeValidator) BindModel(userID uint) error {
	v.Model.UserID = userID
	v.Model.Name = v.Recipe.Name
	v.Model.Description = v.Recipe.Description
	v.Model.PrepTime = v.Recipe.PrepTime
	v.Model.Servings = v.Recipe.Servings
	v.Model.Image = v.Recipe.Image
	if err := v.Model.setTags(v.Recipe.Tags); err != nil {
		return err
	}
	if err := v.Model.setSteps(v.Recipe.Steps); err != nil {
		return err
	}
	if err := v.Model.setIngredientGroups(v.Recipe.IngredientGroups); err != nil {
		return err
	}
	if err := v.Model.setDependencies(v.Recipe.DependentRecipes); err != nil {
		return err
	}
	return nil
}
