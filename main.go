package main

import (
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/router"
	"github.com/anthonyhawkins/savorbook/users"
	"github.com/gofiber/fiber/v2"
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
	router.SetupRoutes(app)
	app.Listen(":3000")
}
