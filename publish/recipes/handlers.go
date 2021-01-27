package recipes

import (
	"errors"
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

func GetRecipes(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	userId := middleware.AuthedUserId(c.Locals("user"))
	searchString := strings.ToLower(c.Query("name"))

	db := database.GetDB()
	var recipes []Recipe
	selects := []string{"id", "user_id", "name", "image", "description"}
	if len(searchString) > 0 {
		_ = db.Select(selects).Where(map[string]interface{}{
			"user_id": userId,
		}).Where("LOWER(name) LIKE ?", "%"+searchString+"%").Find(&recipes)
	} else {
		_ = db.Select(selects).Where(map[string]interface{}{
			"user_id": userId,
		}).Find(&recipes)
	}

	response.Success = true
	response.Message = "All Recipes"
	response.Data = recipes

	return c.JSON(response)

}

func GetRecipe(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	userId := middleware.AuthedUserId(c.Locals("user"))
	recipeId := c.Params("id")

	db := database.GetDB()

	var recipe = new(Recipe)
	result := db.Where(map[string]interface{}{
		"id":      recipeId,
		"user_id": userId,
	}).Preload("Steps.StepImages").Preload("IngredientGroups.Ingredients").First(&recipe)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Message = "Recipe Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if result.Error != nil {
		response.Message = "Unable to Retrieve Recipe"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	db.Model(&Recipe{}).Select("recipes.name as recipe_name, recipe_dependencies.parent_recipe, recipe_dependencies.dependent_recipe, recipe_dependencies.qty").Joins("left join recipe_dependencies on recipe_dependencies.dependent_recipe = recipes.id").Where(map[string]interface{}{
		"recipe_dependencies.parent_recipe": recipeId,
	}).Find(&recipe.DependentRecipes)

	response.Success = true
	response.Data = recipe
	return c.JSON(response)

}

func DeleteRecipe(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	userId := middleware.AuthedUserId(c.Locals("user"))
	recipeId := c.Params("id")

	db := database.GetDB()

	result := db.Where(map[string]interface{}{
		"id":      recipeId,
		"user_id": userId,
	}).Preload("Steps").Preload("StepImages").Preload("IngredientGroups.Ingredients").Delete(&Recipe{})

	if result.Error != nil {
		response.Message = "Unable to Delete Recipe"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response.Success = true
	return c.JSON(response)

}

func UpdateRecipe(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	recipe := new(Recipe)
	err := c.BodyParser(recipe)

	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	errs := ValidateRecipe(*recipe)
	if errs != nil {
		response.Errors = errs
		return c.JSON(response)
	}

	userId := middleware.AuthedUserId(c.Locals("user"))
	recipeId, _ := strconv.ParseUint(c.Params("id"), 10, 64)

	/**
	Ensure the JSON id is present and matches the URL param for recipe id
	*/
	if recipeId != uint64(recipe.ID) {
		response.Message = "Invalid Format"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	db := database.GetDB()

	var existingRecipe = new(Recipe)
	result := db.Where(map[string]interface{}{
		"id":      recipeId,
		"user_id": userId,
	}).Preload("Steps.StepImages").Preload("IngredientGroups.Ingredients").First(&existingRecipe)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Message = "Recipe Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	/**
	When a user updates a list items which need to be deleted are not included in the list
	but we need to ensure they are not orphaned. For now delete the existing list and replace
	with the desired list.

	Alternatives would be when user deletes something in the UI, keep track of "ids_to_delete"
	in another list and send that as well or as a different call. Then delete the ids here.
	*/

	tx := db.Begin()
	tx.Where(map[string]interface{}{"recipe_id": recipeId}).Preload("Ingredients").Select(clause.Associations).Delete(IngredientGroup{})
	tx.Where(map[string]interface{}{"recipe_id": recipeId}).Select(clause.Associations).Delete(Step{})
	tx.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Steps.ID").Updates(&recipe)

	if tx.Save(&recipe); tx.Error != nil {
		tx.Rollback()
		response.Message = "Unable to Update Recipe"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	//Delete Dependencies and Re-add them
	tx.Where(map[string]interface{}{"parent_recipe": recipeId}).Delete(RecipeDependency{})
	if len(recipe.DependentRecipes) > 0 {
		var idsToCheck []uint
		for index, dependency := range recipe.DependentRecipes {
			/**
			Assume the dependent recipe exists, collect its ID
			in the meantime then check against DB when done, once.
			*/
			idsToCheck = append(idsToCheck, dependency.DependentRecipe)
			recipe.DependentRecipes[index].ParentRecipe = recipe.ID
		}

		//Check the DB to ensure all the dependent recipes exists and the user owns them
		var recipes []Recipe
		existResult := db.Where(map[string]interface{}{
			"user_id": userId,
		}).Find(&recipes, idsToCheck)

		if existResult.RowsAffected != int64(len(idsToCheck)) {
			tx.Rollback()
			response.Message = "One or More Dependent Recipes do not exists."
			response.Errors = append(response.Errors, response.Message)
			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		//Create the Recipe Dependency
		createResult := tx.Create(recipe.DependentRecipes)
		if createResult.RowsAffected == 0 {
			tx.Rollback()
			response.Message = "Recipe Update Failed"
			//TODO - Should error be bubbled up DB error to client?
			response.Errors = append(response.Errors, result.Error.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	tx.Commit()
	response.Success = true
	response.Message = "Recipe has been updated"
	response.Data = recipe
	return c.JSON(response)

}

func CreateRecipe(c *fiber.Ctx) error {

	response := new(responses.StandardResponse)
	response.Success = false

	recipe := new(Recipe)
	err := c.BodyParser(recipe)

	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	errs := ValidateRecipe(*recipe)
	if errs != nil {
		response.Errors = errs
		return c.JSON(response)
	}

	userId := middleware.AuthedUserId(c.Locals("user"))
	recipe.UserID = userId

	db := database.GetDB()
	tx := db.Begin()
	//TODO - Review if anything else should be omitted
	result := tx.Omit("Steps.ID").Create(&recipe)

	if result.RowsAffected == 0 {
		tx.Rollback()
		response.Message = "Recipe Creation Failed"
		//TODO - Should error be bubbled up DB error to client?
		response.Errors = append(response.Errors, result.Error.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	//Add Dependencies
	if len(recipe.DependentRecipes) > 0 {
		var idsToCheck []uint
		for index, dependency := range recipe.DependentRecipes {
			/**
			Assume the dependent recipe exists, collect its ID
			in the meantime then check against DB when done, once.
			*/
			idsToCheck = append(idsToCheck, dependency.DependentRecipe)
			recipe.DependentRecipes[index].ParentRecipe = recipe.ID
		}

		//Check the DB to ensure all the dependent recipes exists and the user owns them
		var recipes []Recipe
		existResult := db.Where(map[string]interface{}{
			"user_id": userId,
		}).Find(&recipes, idsToCheck)

		if existResult.RowsAffected != int64(len(idsToCheck)) {
			tx.Rollback()
			response.Message = "One or More Dependent Recipes do not exists."
			response.Errors = append(response.Errors, response.Message)
			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		//Create the Recipe Dependency
		createResult := tx.Create(recipe.DependentRecipes)
		if createResult.RowsAffected == 0 {
			tx.Rollback()
			response.Message = "Recipe Creation Failed"
			//TODO - Should error be bubbled up DB error to client?
			response.Errors = append(response.Errors, result.Error.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	tx.Commit()
	//Respond with Success
	response.Success = true
	response.Data = recipe
	return c.Status(fiber.StatusCreated).JSON(response)

}
