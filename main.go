package main

import (
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/router"
	"github.com/anthonyhawkins/savorbook/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&users.User{})
}

func main() {

	db := database.Init()
	Migrate(db)

	sqlDB := database.GetSqlDB(db)
	defer sqlDB.Close()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)
	app.Listen(":3000")
}
