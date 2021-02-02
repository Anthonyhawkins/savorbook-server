package users

import (
	"github.com/anthonyhawkins/savorbook/database"
)

type UserModel struct {
	database.BaseModel
	Username     string `gorm:"unique_index"`
	Email        string `gorm:"unique_index"`
	DisplayName  string
	Bio          string
	Salt         string
	PasswordHash string
	Status       string
}

func (model *UserModel) Exists() bool {
	db := database.GetDB()
	var existingUsers []UserModel
	db.Where("username = ?", model.Username).Or("email = ?", model.Email).Find(&existingUsers)
	if len(existingUsers) > 0 {
		return true
	}
	return false
}

func (model *UserModel) Create() {
	db := database.GetDB()
	db.Create(&model)
}

func (model *UserModel) Get() {
	db := database.GetDB()
	query := map[string]interface{}{"email": model.Email}
	db.Where(query).Find(&model)
}

func FindOne(userID uint) (*UserModel, error) {
	// Retrieve Existing User and ensure password matches
	db := database.GetDB()
	var user = new(UserModel)
	result := db.First(&user, userID)
	return user, result.Error
}
