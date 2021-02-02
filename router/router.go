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

	//app.Post("/test", middleware.Protected(), recipes.CreateRecipe)

	api := app.Group("/api")

	//Auth
	auth := api.Group("/auth")
	auth.Post("/register", users.UserCreate)
	auth.Post("/login", users.UserLogin)
	auth.Get("/account", middleware.Protected(), users.GetAccount)

	// Publishing
	publish := api.Group("/publish")
	publish.Post("/recipes", middleware.Protected(), recipes.RecipeCreate)
	publish.Get("/recipes", middleware.Protected(), recipes.RecipeList)
	publish.Get("/recipes/:id", middleware.Protected(), recipes.RecipeGet)
	publish.Put("/recipes/:id", middleware.Protected(), recipes.RecipeUpdate)
	publish.Delete("/recipes/:id", middleware.Protected(), recipes.RecipeDelete)

	//library := api.Group("/library")
	//store := api.Group("/store")

	api.Post("/images", middleware.Protected(), images.UploadImage)

}
