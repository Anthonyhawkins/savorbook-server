package recipes

import (
	"errors"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

func RecipeCreate(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	userID := middleware.AuthedUserId(c.Locals("user"))

	recipeValidator := NewRecipeValidator()
	err := c.BodyParser(recipeValidator)
	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	errs, err := recipeValidator.Validate()
	if err != nil {
		response.Message = "Validation Errors"
		response.Errors = errs
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if err := recipeValidator.BindModel(userID); err != nil {
		response.Message = "Unable to Create Recipe"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if err := SaveRecipe(&recipeValidator.Model); err != nil {
		response.Message = "Unable to Create Recipe"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	var recipeResponse RecipeResponse
	recipeResponse.SerializeRecipe(&recipeValidator.Model)

	//Respond with Success
	response.Success = true
	response.Data = recipeResponse
	return c.Status(fiber.StatusCreated).JSON(response)
}

func RecipeGet(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	recipeID := c.Params("id")
	userID := middleware.AuthedUserId(c.Locals("user"))

	model, err := GetRecipeFull(recipeID, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Message = "Recipe Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err != nil {
		response.Message = "Unable to Retrieve Recipe"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var recipeResponse RecipeResponse
	recipeResponse.SerializeRecipe(&model)

	//Respond with Success
	response.Success = true
	response.Data = recipeResponse
	return c.JSON(response)

}

func TagList(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false
	userID := middleware.AuthedUserId(c.Locals("user"))

	tags, _ := GetTags(userID)
	tagResponse := SerializeTags(tags)

	response.Success = true
	response.Data = tagResponse
	response.Errors = append(response.Errors, response.Message)
	return c.JSON(response)

}

func RecipeList(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false
	userID := middleware.AuthedUserId(c.Locals("user"))
	byName := strings.ToLower(c.Query("name"))
	byTags := strings.ToLower(c.Query("tags"))
	var recipeList []RecipeResponse

	var recipes []RecipeModel
	var err error

	if len(byName) > 0 {
		recipes, err = FindRecipesByName(userID, byName)
	} else if len(byTags) > 0 {
		recipes, err = FindRecipesByTags(userID, byTags)
	} else {
		recipes, err = GetRecipes(userID)
	}

	if err != nil {
		response.Success = true
		response.Data = recipeList
		response.Message = "No Recipes Found"
		response.Errors = append(response.Errors, response.Message)
		return c.JSON(response)
	}

	for _, recipe := range recipes {
		var recipeResponse RecipeResponse
		recipeResponse.SerializeRecipe(&recipe)
		recipeList = append(recipeList, recipeResponse)
	}

	//Respond with Success
	response.Success = true
	response.Data = recipeList
	return c.JSON(response)

}

func RecipeUpdate(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	recipeId := c.Params("id")
	userId := middleware.AuthedUserId(c.Locals("user"))
	existingRecipe, err := GetRecipe(recipeId, userId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Message = "Recipe Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err != nil {
		response.Message = "Unable to Retrieve Recipe"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	recipeValidator := NewRecipeValidator()
	if err := c.BodyParser(recipeValidator); err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	validationErrors, err := recipeValidator.Validate()
	if err != nil {
		response.Message = "Validation Errors"
		response.Errors = validationErrors
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if err := recipeValidator.BindModel(userId); err != nil {
		response.Message = "Unable to Update Recipe"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	recipeValidator.Model.ID = existingRecipe.ID

	if err := recipeValidator.Model.Update(); err != nil {
		response.Message = "Unable to Update Recipe"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	var recipeResponse RecipeResponse
	recipeResponse.SerializeRecipe(&recipeValidator.Model)

	//Respond with Success
	response.Success = true
	response.Data = recipeResponse
	return c.JSON(response)
}

func RecipeDelete(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	recipeId := c.Params("id")
	userId := middleware.AuthedUserId(c.Locals("user"))

	parentRecipes, err := DeleteRecipe(recipeId, userId)

	if err != nil {
		response.Message = "Unable to Delete Recipe."
		response.Errors = append(response.Errors, err.Error())
		response.Data = SerializeParentRecipes(parentRecipes)
		return c.JSON(response)
	}

	//Respond with Success
	response.Success = true
	return c.JSON(response)

}
