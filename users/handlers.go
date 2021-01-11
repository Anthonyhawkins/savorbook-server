package users

import (
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

func CreateUser(c *fiber.Ctx) error {

	response := new(responses.StandardResponse)
	response.Success = false

	registration := new(RegisterForm)
	err := c.BodyParser(registration)
	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	errors := ValidateRegistration(*registration)
	if errors != nil {
		return c.JSON(errors)
	}

	//Check to see if username and email is already in use.
	db := database.GetDB()

	var existingUsers []User
	db.Where("username = ?", registration.Username).Or("email = ?", registration.Email).Find(&existingUsers)

	if len(existingUsers) > 0 {
		response.Message = "Username and or Email already in use."
		for _, existingUser := range existingUsers {
			if existingUser.Email == registration.Email {
				response.Errors = append(response.Errors, "An account with this email already exists.")
			}
			if existingUser.Username == registration.Username {
				response.Errors = append(response.Errors, "This username has already been taken.")
			}
		}
		return c.JSON(response)
	}

	// create the new user
	user := new(User)
	user.PasswordHash = setPassword(registration.Password)
	user.Username = registration.Username
	user.DisplayName = registration.DisplayName
	user.Email = registration.Email
	db.Create(&user)

	//generate JWT Token
	signedToken, err := middleware.SetToken(user.Username, user.Email, user.ID)
	if err != nil {
		response.Message = "Login Error"
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response.Success = true
	response.Message = "Registration Successful"

	response.Data = struct {
		AccessToken string `json:"accessToken"`
		UserId      uint   `json:"userId"`
		DisplayName string `json:"displayName"`
	}{
		AccessToken: signedToken,
		UserId:      user.ID,
		DisplayName: user.DisplayName,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func LogInUser(c *fiber.Ctx) error {

	response := new(responses.StandardResponse)
	response.Success = false

	//Parse Login Form
	login := new(LoginForm)
	err := c.BodyParser(login)
	if err != nil {
		response.Message = "Invalid JSON"
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Retrieve Existing User and ensure password matches
	db := database.GetDB()
	query := map[string]interface{}{"username": login.Username}
	var user = new(User)
	db.Where(query).Find(&user)

	if !checkPassword(login.Password, user.PasswordHash) {
		response.Message = "Invalid Login"
		response.Errors = append(response.Errors, "Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	//generate JWT Token and return with user data
	signedToken, err := middleware.SetToken(user.Username, user.Email, user.ID)
	if err != nil {
		response.Message = "Login Error"
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response.Success = true
	response.Message = "Login Successful"

	response.Data = struct {
		AccessToken string `json:"accessToken"`
		UserId      uint   `json:"userId"`
		DisplayName string `json:"displayName"`
	}{
		AccessToken: signedToken,
		UserId:      user.ID,
		DisplayName: user.DisplayName,
	}

	return c.JSON(response)

}
