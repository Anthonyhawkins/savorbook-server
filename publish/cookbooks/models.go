package cookbooks

import (
	"errors"
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/publish/recipes"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type CookbookModel struct {
	gorm.Model
	UserID   uint
	Title    string
	SubTitle string
	Image    string
	Blurb    string
	Sections []SectionModel `gorm:"foreignKey:CookbookID;constraint:OnDelete:CASCADE"`
}

type SectionModel struct {
	gorm.Model
	UserID     uint
	Name       string
	Overview   string
	Recipes    pq.Int64Array `gorm:"type:integer[]"`
	CookbookID uint
}

//TODO - There could be a use for a "page type" abstraction to help with keeping the order in which
// a recipe appears in a section or cookbook
// A page could also be an interface in which multiple entities could satisfy, i.e. a Recipe, Guide, etc..
/*type PageModel struct {
	gorm.Model
	UserID uint
	PageType string
	Location string
	ItemID uint
}*/

func (model *CookbookModel) setSections(userID uint, sectionValidators []SectionValidator) error {
	var sections []SectionModel
	for _, sectionValidator := range sectionValidators {
		var section SectionModel
		section.UserID = userID
		section.Name = sectionValidator.Name
		section.Overview = sectionValidator.Overview

		existingRecipes, err := recipes.GetRecipesByIDs(userID, sectionValidator.Recipes)
		if err != nil {
			return err
		}

		if len(existingRecipes) != len(sectionValidator.Recipes) {
			return errors.New("one or more recipes do not exist")
		}

		var recipeIDs []int64
		for _, recipeID := range sectionValidator.Recipes {
			recipeIDs = append(recipeIDs, int64(recipeID))
		}

		section.Recipes = recipeIDs
		sections = append(sections, section)
	}
	model.Sections = sections
	return nil
}

func DeleteCookbook(cookbookID string, userID uint) error {
	db := database.GetDB()

	result := db.Where(map[string]interface{}{
		"id":      cookbookID,
		"user_id": userID,
	}).Delete(&CookbookModel{})

	if result.Error != nil {
		return errors.New("unable to delete cookbook")
	}

	return nil
}

func CreateCookbook(cookbook *CookbookModel) error {
	db := database.GetDB()

	err := db.Create(cookbook).Error
	return err
}

func GetCookbook(cookbookID string, userID uint) (CookbookModel, error) {
	db := database.GetDB()
	var model CookbookModel

	result := db.Where(map[string]interface{}{
		"id":      cookbookID,
		"user_id": userID,
	}).Preload("Sections").First(&model)

	return model, result.Error

}

func GetSectionRecipes(sectionID string, userID uint) ([]recipes.RecipeModel, error) {

	recipesList := make([]recipes.RecipeModel, 0)

	db := database.GetDB()
	var section SectionModel

	result := db.Where(map[string]interface{}{
		"id":      sectionID,
		"user_id": userID,
	}).First(&section)

	if result.Error != nil {
		return recipesList, result.Error
	}
	/**
	Why a Loop and not an IDs IN Query? We need to keep the Recipes in the same order
	as their IDs as this is the order they are displayed in the section and ultimately in
	the cookbook.
	*/
	for _, recipeID := range section.Recipes {
		recipe, err := recipes.GetRecipe(strconv.FormatInt(recipeID, 10), userID)
		if err == nil {
			recipesList = append(recipesList, recipe)
		}
	}

	return recipesList, result.Error
}

func (model *CookbookModel) Update() error {
	db := database.GetDB()
	tx := db.Begin()
	tx.Where(map[string]interface{}{"cookbook_id": model.ID}).Select(clause.Associations).Delete(&SectionModel{})
	if tx.Save(&model); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func GetCookbooks(userID uint, pageNum string, pageSize string) ([]CookbookModel, error) {

	db := database.GetDB()
	var cookbooks []CookbookModel

	result := db.Scopes(database.Paginate(pageNum, pageSize)).Where(map[string]interface{}{
		"user_id": userID,
	}).Preload("Sections").Find(&cookbooks)

	return cookbooks, result.Error
}
