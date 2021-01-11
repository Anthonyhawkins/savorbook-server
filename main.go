package main

import (
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/handlers"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.UserModel{})
}

func main() {

	db := database.Init()
	Migrate(db)

	sqlDB := database.GetSqlDB(db)
	defer sqlDB.Close()

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/secret", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.SendString("Super Secret Page")
	})

	app.Post("/users", handlers.CreateUser)
	app.Post("/login", handlers.LogInUser)

	app.Listen(":3000")
}
