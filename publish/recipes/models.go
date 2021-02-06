package recipes

import (
	"errors"
	"fmt"
	"github.com/anthonyhawkins/savorbook/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type RecipeModel struct {
	gorm.Model
	UserID           uint
	Name             string
	Image            string
	Description      string
	PrepTime         string
	Servings         string
	Tags             []TagModel              `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	DependentRecipes []RecipeDependencyModel `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	ParentRecipes    []RecipeDependencyModel `gorm:"-"`
	IngredientGroups []IngredientGroupModel  `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Steps            []StepModel             `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
}

type TagModel struct {
	gorm.Model
	RecipeID uint
	UserID   uint
	Tag      string
}

type IngredientGroupModel struct {
	gorm.Model
	GroupName   string
	Ingredients []IngredientModel `gorm:"foreignKey:IngredientGroupID;constraint:OnDelete:CASCADE"`
	RecipeID    uint
}

type IngredientModel struct {
	gorm.Model
	Name              string
	Qty               string
	Unit              string
	IngredientGroupID uint
}

type StepModel struct {
	gorm.Model
	Type       string
	Text       string
	StepImages []StepImageModel `gorm:"foreignKey:StepID;constraint:OnDelete:CASCADE"`
	RecipeID   uint
}

type StepImageModel struct {
	gorm.Model
	Image  string
	Text   string
	StepID uint
}

type RecipeDependencyModel struct {
	gorm.Model
	RecipeID        uint
	DependentRecipe uint
	RecipeName      string
	Qty             string
}

func SaveRecipe(recipe *RecipeModel) error {
	db := database.GetDB()

	if err := recipe.CheckDependencies(db); err != nil {
		return err
	}

	err := db.Save(recipe).Error
	return err
}

func GetRecipeParents(recipeID string) ([]RecipeDependencyModel, error) {

	db := database.GetDB()
	parentRecipes := make([]RecipeDependencyModel, 0)
	result := db.Model(&RecipeModel{}).Select(
		`recipe_models.id,
        recipe_models.name as recipe_name, 
		recipe_dependency_models.recipe_id, 
        recipe_dependency_models.dependent_recipe, 
        recipe_dependency_models.qty`,
	).Joins(
		`left join recipe_dependency_models 
        on recipe_dependency_models.recipe_id = recipe_models.id`,
	).Where(map[string]interface{}{
		"recipe_dependency_models.dependent_recipe": recipeID,
		"recipe_dependency_models.deleted_at":       nil,
	}).Find(&parentRecipes)

	return parentRecipes, result.Error

}

func DeleteRecipe(recipeID string, userID uint) ([]RecipeDependencyModel, error) {
	db := database.GetDB()

	parentRecipes, _ := GetRecipeParents(recipeID)

	if len(parentRecipes) > 0 {
		return parentRecipes, errors.New("this recipe is listed as a dependent for another recipe")
	}

	result := db.Where(map[string]interface{}{
		"id":      recipeID,
		"user_id": userID,
	}).Delete(&RecipeModel{})

	if result.Error != nil {
		return parentRecipes, errors.New("unable to delete recipe")
	}

	return parentRecipes, nil
}

func (model *RecipeModel) Update() error {
	db := database.GetDB()

	if err := model.CheckDependencies(db); err != nil {
		return err
	}

	/**
	When a user updates a list, items which need to be deleted are not included in the list
	but we need to ensure they are not orphaned. For now delete the existing list and replace
	with the desired list.
	Alternatives would be when user deletes something in the UI, keep track of "ids_to_delete"
	in another list and send that as well or as a different call. Then delete the ids here.
	*/
	tx := db.Begin()
	//tx.Model(&model).Association("Tags").Replace(model.Tags)
	tx.Where(map[string]interface{}{"recipe_id": model.ID}).Select(clause.Associations).Delete(&TagModel{})
	tx.Where(map[string]interface{}{"recipe_id": model.ID}).Select(clause.Associations).Delete(&IngredientGroupModel{})
	tx.Where(map[string]interface{}{"recipe_id": model.ID}).Select(clause.Associations).Delete(&StepModel{})
	tx.Where(map[string]interface{}{"recipe_id": model.ID}).Select(clause.Associations).Delete(&RecipeDependencyModel{})

	if tx.Save(&model); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (model *RecipeModel) setTags(tags []string) error {
	var tagList []TagModel

	for _, tag := range tags {
		var tagModel TagModel
		tagModel.UserID = model.UserID
		tagModel.Tag = tag
		tagList = append(tagList, tagModel)
	}

	model.Tags = tagList
	return nil

}

func (model *RecipeModel) CheckDependencies(db *gorm.DB) error {
	/**
	Assume the dependent recipe exists, collect their IDs
	in the meantime then check against DB when done, once.
	*/

	var idsToCheck []uint
	for _, dependency := range model.DependentRecipes {
		idsToCheck = append(idsToCheck, dependency.DependentRecipe)
	}

	if len(idsToCheck) == 0 {
		return nil
	}

	var recipes []RecipeModel
	existResult := db.Where(map[string]interface{}{
		"user_id": model.UserID,
	}).Find(&recipes, idsToCheck)

	if existResult.RowsAffected != int64(len(idsToCheck)) {
		return errors.New("one or more dependent recipes does not exist")
	}

	return nil
}

func (model *RecipeModel) setDependencies(dependents []RecipeDependencyValidator) error {
	var dependencies []RecipeDependencyModel
	for _, dependent := range dependents {
		var recipeDependency RecipeDependencyModel
		recipeDependency.DependentRecipe = dependent.DependentRecipe
		recipeDependency.Qty = dependent.Qty
		dependencies = append(dependencies, recipeDependency)
	}
	model.DependentRecipes = dependencies
	return nil
}

func (model *RecipeModel) setIngredientGroups(ingredientGroupValidators []IngredientGroupValidator) error {
	var groups []IngredientGroupModel
	for _, groupValidator := range ingredientGroupValidators {
		var group IngredientGroupModel
		group.GroupName = groupValidator.GroupName
		if err := group.setIngredients(groupValidator.Ingredients); err != nil {
			return err
		}
		groups = append(groups, group)
	}
	model.IngredientGroups = groups
	return nil
}

func (model *IngredientGroupModel) setIngredients(ingredientValidators []IngredientValidator) error {
	var ingredients []IngredientModel
	for _, ingredientValidator := range ingredientValidators {
		var ingredient IngredientModel
		ingredient.Name = ingredientValidator.Name
		ingredient.Qty = ingredientValidator.Qty
		ingredient.Unit = ingredientValidator.Unit
		ingredients = append(ingredients, ingredient)
	}
	model.Ingredients = ingredients
	return nil
}

func (model *RecipeModel) setSteps(stepValidators []StepValidator) error {
	var steps []StepModel
	for _, stepValidator := range stepValidators {
		var step StepModel
		step.Type = stepValidator.Type
		step.Text = stepValidator.Text
		if err := step.setStepImages(stepValidator.StepImages); err != nil {
			return err
		}
		steps = append(steps, step)
	}
	model.Steps = steps
	return nil
}

func (model *StepModel) setStepImages(stepImageValidators []StepImageValidator) error {
	var images []StepImageModel
	for _, stepImageValidator := range stepImageValidators {
		var stepImage StepImageModel
		stepImage.Text = stepImageValidator.Text
		stepImage.Image = stepImageValidator.Image
		images = append(images, stepImage)
	}
	model.StepImages = images
	return nil
}

/**
GET RECIPE
*/
func GetRecipe(recipeID string, userID uint) (RecipeModel, error) {
	db := database.GetDB()
	var model RecipeModel

	result := db.Where(map[string]interface{}{
		"id":      recipeID,
		"user_id": userID,
	}).Preload("Tags").Find(&model)

	return model, result.Error
}

func GetRecipeFull(recipeID string, userID uint) (RecipeModel, error) {
	db := database.GetDB()
	var model RecipeModel

	result := db.Where(map[string]interface{}{
		"id":      recipeID,
		"user_id": userID,
	}).Preload("Tags").Preload("Steps.StepImages").Preload("IngredientGroups.Ingredients").First(&model)

	if result.Error != nil {
		return model, result.Error
	}

	dependencyResult := db.Model(&RecipeModel{}).Select(
		`recipe_models.id,
        recipe_models.name as recipe_name, 
		recipe_dependency_models.recipe_id, 
        recipe_dependency_models.dependent_recipe, 
        recipe_dependency_models.qty`,
	).Joins(
		`left join recipe_dependency_models 
        on recipe_dependency_models.dependent_recipe = recipe_models.id`,
	).Where(map[string]interface{}{
		"recipe_dependency_models.recipe_id":  recipeID,
		"recipe_dependency_models.deleted_at": nil,
	}).Find(&model.DependentRecipes)

	if dependencyResult.Error != nil {
		return model, result.Error
	}

	parentRecipes, err := GetRecipeParents(fmt.Sprint(model.ID))
	model.ParentRecipes = parentRecipes
	return model, err
}

func GetRecipes(userID uint) ([]RecipeModel, error) {

	db := database.GetDB()
	var recipes []RecipeModel

	selects := []string{"id", "user_id", "name", "image", "description", "prep_time", "servings"}
	result := db.Select(selects).Where(map[string]interface{}{
		"user_id": userID,
	}).Preload("Tags").Find(&recipes)

	return recipes, result.Error
}

func FindRecipesByName(userID uint, searchString string) ([]RecipeModel, error) {

	db := database.GetDB()
	var recipes []RecipeModel

	selects := []string{"id", "user_id", "name", "image", "description", "prep_time", "servings"}
	result := db.Select(selects).Where(map[string]interface{}{
		"user_id": userID,
	}).Where("LOWER(name) LIKE ?", "%"+searchString+"%").Find(&recipes)

	return recipes, result.Error
}

func FindRecipesByTags(userID uint, searchString string) ([]RecipeModel, error) {

	tags := strings.Split(strings.ToLower(searchString), ",")

	db := database.GetDB()
	var recipes []RecipeModel

	result := db.Model(&RecipeModel{}).Distinct().Preload("Tags").Joins(
		`left join tag_models 
        on tag_models.recipe_id = recipe_models.id`,
	).Where(map[string]interface{}{
		"tag_models.user_id":    userID,
		"tag_models.deleted_at": nil,
	}).Where("tag_models.tag IN ?", tags).Find(&recipes)

	if result.Error != nil {
		return recipes, result.Error
	}

	return recipes, result.Error
}

func GetTags(userID uint) ([]TagModel, error) {

	db := database.GetDB()
	var tags []TagModel
	result := db.Where(map[string]interface{}{
		"user_id": userID,
	}).Distinct("tag").Order("tag").Find(&tags)

	return tags, result.Error
}
