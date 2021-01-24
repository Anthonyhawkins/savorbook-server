package router

import (
	"github.com/anthonyhawkins/savorbook/images"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/publish/recipes"
	"github.com/anthonyhawkins/savorbook/users"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	//Home
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/test", middleware.Protected(), recipes.CreateRecipe)

	api := app.Group("/api")

	//Auth
	auth := api.Group("/auth")
	auth.Post("/register", users.CreateUser)
	auth.Post("/login", users.LogInUser)
	auth.Get("/account", middleware.Protected(), users.GetAccount)

	// Publishing
	publish := api.Group("/publish")
	publish.Post("/recipes", middleware.Protected(), recipes.CreateRecipe)
	publish.Get("/recipes", middleware.Protected(), recipes.GetRecipes)
	publish.Get("/recipes/:id", middleware.Protected(), recipes.GetRecipe)
	publish.Put("/recipes/:id", middleware.Protected(), recipes.UpdateRecipe)
	publish.Delete("/recipes/:id", middleware.Protected(), recipes.DeleteRecipe)

	//library := api.Group("/library")
	//store := api.Group("/store")

	api.Post("/images", middleware.Protected(), images.UploadImage)

}
