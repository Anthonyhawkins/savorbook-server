package recipes

import (
	"github.com/anthonyhawkins/savorbook/database"
)

type Recipe struct {
	database.BaseModel
	UserID           uint              `json:"userId"`
	Name             string            `json:"name"              validate:"required,max=75"`
	Description      string            `json:"description"       validate:"required,max=2600"`
	IngredientGroups []IngredientGroup `json:"ingredientGroups"  validate:"required,dive"       gorm:"constraint:OnDelete:CASCADE"`
	Steps            []Step            `json:"steps"             validate:"required,dive"       gorm:"constraint:OnDelete:CASCADE"            `
}

type IngredientGroup struct {
	database.BaseModel
	GroupName   string       `json:"groupName"   validate:"max=50"`
	Ingredients []Ingredient `json:"ingredients" validate:"required,dive" gorm:"constraint:OnDelete:CASCADE" `
	RecipeID    uint         `json:"recipeId"`
}

type Ingredient struct {
	database.BaseModel
	Name              string `json:"name" validate:"required,max=32"`
	Qty               string `json:"qty"  validate:"omitempty,numeric,max=3"`
	Unit              string `json:"unit" validate:"omitempty,max=12"`
	IngredientGroupID uint   `json:"ingredientGroupId"`
}

type Step struct {
	database.BaseModel
	Type        string `json:"type"        validate:"oneof=text tipText"`
	Text        string `json:"text"        validate:"max=865"`
	TextRight   string `json:"textRight"   validate:"max=432"`
	ImageRight  string `json:"imageRight"  validate:"omitempty,url"`
	TextCenter  string `json:"textCenter"  validate:"max=432"`
	ImageCenter string `json:"imageCenter" validate:"omitempty,url"`
	TextLeft    string `json:"textLeft"    validate:"max=432"`
	ImageLeft   string `json:"imageLeft"   validate:"omitempty,url"`
	RecipeID    uint   `json:"recipeId"`
}
