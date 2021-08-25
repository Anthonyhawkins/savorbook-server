package cookbooks

import (
	"errors"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/publish/recipes"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

func CookbookCreate(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	userID := middleware.AuthedUserId(c.Locals("user"))

	cookbookValidator := NewCookbookValidator()
	err := c.BodyParser(cookbookValidator)
	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	errs, err := cookbookValidator.Validate()
	if err != nil {
		response.Message = "Validation Errors"
		response.Errors = errs
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if err := cookbookValidator.BindModel(userID); err != nil {
		response.Message = "Unable to Create Cookbook"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if err := CreateCookbook(&cookbookValidator.Model); err != nil {
		response.Message = "Unable to Create Cookbook"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	var cookbookResponse CookbookResponse
	cookbookResponse.SerializeCookbook(&cookbookValidator.Model)

	//Respond with Success
	response.Success = true
	response.Data = cookbookResponse
	return c.Status(fiber.StatusCreated).JSON(response)
}

func CookbookGet(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	cookbookID := c.Params("id")
	userID := middleware.AuthedUserId(c.Locals("user"))

	model, err := GetCookbook(cookbookID, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Message = "Cookbook Not Fount"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err != nil {
		response.Message = "Unable to Retrieve Cookbook"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var cookbookResponse CookbookResponse
	cookbookResponse.SerializeCookbook(&model)
	response.Success = true
	response.Data = cookbookResponse
	return c.JSON(response)

}

func SectionRecipesGet(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false
	sectionID := c.Params("id")
	userID := middleware.AuthedUserId(c.Locals("user"))

	recipeList := make([]recipes.RecipeResponse, 0)

	recipeModels, err := GetSectionRecipes(sectionID, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Message = "Section Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err != nil {
		response.Message = "Unable to Retrieve Section"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	for _, recipeModel := range recipeModels {
		var recipeResponse recipes.RecipeResponse
		recipeResponse.SerializeRecipe(&recipeModel)
		recipeList = append(recipeList, recipeResponse)
	}

	//Respond with Success
	response.Success = true
	response.Data = recipeList
	return c.JSON(response)

}

func CookbookUpdate(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	cookbookID := c.Params("id")
	userID := middleware.AuthedUserId(c.Locals("user"))
	existingCookbook, err := GetCookbook(cookbookID, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Message = "Cookbook Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err != nil {
		response.Message = "Unable to Retrieve Cookbook"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	cookbookValidator := NewCookbookValidator()
	if err := c.BodyParser(cookbookValidator); err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	validationErrors, err := cookbookValidator.Validate()
	if err != nil {
		response.Message = "Validation Errors"
		response.Errors = validationErrors
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if err := cookbookValidator.BindModel(userID); err != nil {
		response.Message = "Unable to Update Cookbook"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	cookbookValidator.Model.ID = existingCookbook.ID

	if err := cookbookValidator.Model.Update(); err != nil {
		response.Message = "Unable to Update Cookbook"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	var cookbookResponse CookbookResponse
	cookbookResponse.SerializeCookbook(&cookbookValidator.Model)

	//Respond with Success
	response.Success = true
	response.Data = cookbookResponse
	return c.JSON(response)
}

func CookbookList(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false
	userID := middleware.AuthedUserId(c.Locals("user"))
	pageNum := strings.ToLower(c.Query("page"))
	pageSize := strings.ToLower(c.Query("page_size"))

	cookbookList := make([]CookbookResponse, 0)

	recipes, err := GetCookbooks(userID, pageNum, pageSize)

	if err != nil {
		response.Success = true
		response.Data = cookbookList
		response.Message = "No Cookbooks Found"
		response.Errors = append(response.Errors, response.Message)
		return c.JSON(response)
	}

	for _, recipe := range recipes {
		var cookbookResponse CookbookResponse
		cookbookResponse.SerializeCookbook(&recipe)
		cookbookList = append(cookbookList, cookbookResponse)
	}

	//Respond with Success
	response.Success = true
	response.Data = cookbookList
	return c.JSON(response)
}

func CookbookDelete(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	recipeId := c.Params("id")
	userId := middleware.AuthedUserId(c.Locals("user"))

	err := DeleteCookbook(recipeId, userId)

	if err != nil {
		response.Message = "Unable to Delete Cookbook."
		response.Errors = append(response.Errors, err.Error())
		return c.JSON(response)
	}

	//Respond with Success
	response.Success = true
	return c.JSON(response)

}
