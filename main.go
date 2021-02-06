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

	db.Migrator().DropTable(&recipes.RecipeModel{})
	db.Migrator().DropTable(&recipes.TagModel{})
	type RecipeTags struct {
	}
	db.Migrator().DropTable(&RecipeTags{})
	db.Migrator().DropTable(&recipes.IngredientGroupModel{})
	db.Migrator().DropTable(&recipes.IngredientModel{})
	db.Migrator().DropTable(&recipes.StepModel{})
	db.Migrator().DropTable(&recipes.StepImageModel{})
	db.Migrator().DropTable(&recipes.RecipeDependencyModel{})
	//db.Migrator().DropTable(&users.UserModel{})

	db.AutoMigrate(&users.UserModel{})
	db.AutoMigrate(&recipes.RecipeModel{})
	db.AutoMigrate(&recipes.TagModel{})
	db.AutoMigrate(&recipes.IngredientGroupModel{})
	db.AutoMigrate(&recipes.IngredientModel{})
	db.AutoMigrate(&recipes.StepModel{})
	db.AutoMigrate(&recipes.StepImageModel{})
	db.AutoMigrate(&recipes.RecipeDependencyModel{})
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
