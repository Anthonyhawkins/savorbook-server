package users

import (
	"errors"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setPassword(password string) string {
	bytePassword := []byte(password)
	hash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(hash)
}

func checkPassword(loginPassword string, userPassword string) bool {
	byteLoginPassword := []byte(loginPassword)
	byteUserPassword := []byte(userPassword)
	err := bcrypt.CompareHashAndPassword(byteUserPassword, byteLoginPassword)
	if err != nil {
		return false
	}

	return true
}

func UserCreate(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	registrationValidator := NewRegisterValidator()
	err := c.BodyParser(registrationValidator)
	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	errs, err := registrationValidator.Validate()
	if errs != nil {
		response.Errors = errs
		return c.JSON(response)
	}

	if err := registrationValidator.BindModel(); err != nil {
		response.Message = "Unable to Create User"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	if registrationValidator.Model.Exists() {
		response.Message = "Username and or Email already in use."
		response.Errors = append(response.Errors, response.Message)
		return c.JSON(response)
	}

	registrationValidator.Model.PasswordHash = setPassword(registrationValidator.Registration.Password)
	registrationValidator.Model.Create()

	//generate JWT Token
	signedToken, err := middleware.SetToken(
		registrationValidator.Model.Username,
		registrationValidator.Model.DisplayName,
		registrationValidator.Model.Email,
		registrationValidator.Model.ID,
	)
	if err != nil {
		response.Message = "Login Error"
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var loginResponse LoginResponse
	loginResponse.SerializeLogin(&registrationValidator.Model, signedToken)

	response.Success = true
	response.Message = "Registration Successful"
	response.Data = loginResponse
	return c.Status(fiber.StatusCreated).JSON(response)

}

func UserLogin(c *fiber.Ctx) error {
	response := new(responses.StandardResponse)
	response.Success = false

	loginValidator := NewLoginValidator()
	err := c.BodyParser(loginValidator)
	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	errs, err := loginValidator.Validate()
	if errs != nil {
		response.Errors = errs
		return c.JSON(response)
	}

	if err := loginValidator.BindModel(); err != nil {
		response.Message = "Unable Login User"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}

	loginValidator.Model.Get()
	if !checkPassword(loginValidator.Login.Password, loginValidator.Model.PasswordHash) {
		response.Message = "Invalid Login"
		response.Errors = append(response.Errors, "Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	//generate JWT Token
	signedToken, err := middleware.SetToken(
		loginValidator.Model.Username,
		loginValidator.Model.DisplayName,
		loginValidator.Model.Email,
		loginValidator.Model.ID,
	)
	if err != nil {
		response.Message = "Login Error"
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var loginResponse LoginResponse
	loginResponse.SerializeLogin(&loginValidator.Model, signedToken)

	response.Success = true
	response.Message = "Login Successful"
	response.Data = loginResponse
	return c.Status(fiber.StatusCreated).JSON(response)

}

func GetAccount(c *fiber.Ctx) error {

	response := new(responses.StandardResponse)
	response.Success = false

	userID := middleware.AuthedUserId(c.Locals("user"))

	userModel, err := FindOne(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Message = "Account Not Found"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	var userResponse UserResponse
	userResponse.SerializeUser(userModel)

	response.Success = true
	response.Message = "Account Retrieval Successful"
	response.Data = userResponse
	return c.JSON(response)

}
