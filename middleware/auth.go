package middleware

import (
	"github.com/anthonyhawkins/savorbook/config"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"time"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.Get("SIGNING_SECRET")),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {

	response := responses.StandardResponse{
		Success: false,
		Message: "Bad Request",
	}

	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

func SetToken(userName string, displayName string, email string, userID uint) (string, error) {
	//generate JWT Token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userName
	claims["displayName"] = displayName
	claims["email"] = email
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString([]byte(config.Get("SIGNING_SECRET")))

	return signedToken, err
}
