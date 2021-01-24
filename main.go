package main

import (
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/images"
	"github.com/anthonyhawkins/savorbook/publish/recipes"
	"github.com/anthonyhawkins/savorbook/router"
	"github.com/anthonyhawkins/savorbook/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	db.Migrator().DropTable(&recipes.Recipe{})
	db.Migrator().DropTable(&recipes.IngredientGroup{})
	db.Migrator().DropTable(&recipes.Ingredient{})
	db.Migrator().DropTable(&recipes.Step{})
	db.Migrator().DropTable(&recipes.StepImage{})
	//db.Migrator().DropTable(&users.User{})

	db.AutoMigrate(&users.User{})
	db.AutoMigrate(&recipes.Recipe{})
	db.AutoMigrate(&recipes.IngredientGroup{})
	db.AutoMigrate(&recipes.Ingredient{})
	db.AutoMigrate(&recipes.Step{})
	db.AutoMigrate(&recipes.StepImage{})
	db.AutoMigrate(&images.Image{})
}

func main() {

	db := database.Init()
	Migrate(db)

	sqlDB := database.GetSqlDB(db)
	defer sqlDB.Close()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		BodyLimit:     4194304,
	})
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)
	app.Listen(":3000")
}
