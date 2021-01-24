package recipes

import (
	"github.com/anthonyhawkins/savorbook/database"
)

type Recipe struct {
	database.BaseModel
	UserID           uint              `json:"userId"`
	Name             string            `json:"name"              validate:"required,max=75"`
	Image            string            `json:"image"             validate:"omitempty,url"`
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
	Qty               string `json:"qty"  validate:"omitempty,max=4"`
	Unit              string `json:"unit" validate:"omitempty,max=12"`
	IngredientGroupID uint   `json:"ingredientGroupId"`
}

type Step struct {
	database.BaseModel
	Type       string      `json:"type"        validate:"oneof=text tipText imageLeft imageRight imageDouble imageTriple"`
	Text       string      `json:"text"        validate:"max=865"`
	StepImages []StepImage `json:"images"      validate:"dive" gorm:"constraint:OnDelete:CASCADE" `
	RecipeID   uint        `json:"recipeId"`
}

type StepImage struct {
	database.BaseModel
	Image  string `json:"src"`
	Text   string `json:"text"`
	StepID uint   `json:"stepId"`
}
