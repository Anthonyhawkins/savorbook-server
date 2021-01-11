package handlers

import (
	"github.com/anthonyhawkins/savorbook/config"
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/forms"
	"github.com/anthonyhawkins/savorbook/models"
	"github.com/anthonyhawkins/savorbook/validators"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
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

	registration := new(forms.RegisterForm)
	err := c.BodyParser(registration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := validators.ValidateRegistration(*registration)
	if errors != nil {
		return c.JSON(errors)
	}

	user := new(models.UserModel)
	user.PasswordHash = setPassword(registration.Password)
	user.Username = registration.Username
	user.Email = registration.Email

	db := database.GetDB()
	db.Create(&user)
	return c.JSON(user)
}

func LogInUser(c *fiber.Ctx) error {

	login := new(forms.LoginForm)
	err := c.BodyParser(login)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	db := database.GetDB()

	// Retrieve Existing User and ensure password matches
	query := map[string]interface{}{"username": login.Username}
	var user = new(models.UserModel)
	db.Where(query).Find(&user)

	if !checkPassword(login.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Login",
		})
	}

	//generate JWT Token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["sub"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString([]byte(config.Get("SIGNING_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Login Error",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Login Successful",
		"data":    signedToken,
	})

}
