package router

import (
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/users"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	//Home
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	api := app.Group("/api")

	//Auth
	auth := api.Group("/auth")
	auth.Post("/register", users.CreateUser)
	auth.Post("/login", users.LogInUser)
	auth.Get("/account", middleware.Protected(), users.GetAccount)

	// User
	//publish := api.Group("/publish")
	//library := api.Group("/library")
	//store := api.Group("/store")

	app.Get("/secret", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.SendString("Super Secret Page")
	})

}
