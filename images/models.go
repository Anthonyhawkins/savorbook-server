package images

import (
	"github.com/anthonyhawkins/savorbook/database"
)

type Image struct {
	database.BaseModel
	RecipeID uint   `json:"recipeId"`
	UserID   uint   `json:"userId"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Used     bool   `json:"used"`
}
