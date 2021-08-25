package router

import (
	"github.com/anthonyhawkins/savorbook/images"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/publish/cookbooks"
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
	auth.Put("/account", middleware.Protected(), users.UpdateAccount)
	auth.Put("/account/password", middleware.Protected(), users.UpdatePassword)

	// Publishing
	publish := api.Group("/publish")
	publish.Post("/recipes", middleware.Protected(), recipes.RecipeCreate)
	publish.Get("/recipes", middleware.Protected(), recipes.RecipeList)
	publish.Get("/recipes/tags", middleware.Protected(), recipes.TagList)
	publish.Get("/recipes/:id", middleware.Protected(), recipes.RecipeGet)
	publish.Put("/recipes/:id", middleware.Protected(), recipes.RecipeUpdate)
	publish.Delete("/recipes/:id", middleware.Protected(), recipes.RecipeDelete)

	publish.Post("/cookbooks", middleware.Protected(), cookbooks.CookbookCreate)
	publish.Get("/cookbooks", middleware.Protected(), cookbooks.CookbookList)
	publish.Get("/cookbooks/:id", middleware.Protected(), cookbooks.CookbookGet)
	publish.Put("/cookbooks/:id", middleware.Protected(), cookbooks.CookbookUpdate)
	publish.Delete("/cookbooks/:id", middleware.Protected(), cookbooks.CookbookDelete)
	publish.Get("/sections/:id/recipes", middleware.Protected(), cookbooks.SectionRecipesGet)
	//library := api.Group("/library")
	//store := api.Group("/store")

	api.Post("/images", middleware.Protected(), images.UploadImage)

}
