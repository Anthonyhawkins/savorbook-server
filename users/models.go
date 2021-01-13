package users

import (
	"github.com/anthonyhawkins/savorbook/database"
)

/**
Data to be modeled into the Database
*/

type User struct {
	database.BaseModel
	Username     string `gorm:"column:username;unique_index" json:"username"`
	Email        string `gorm:"column:email;unique_index" json:"email"`
	DisplayName  string `gorm:"column:display_name" json:"displayName"`
	Bio          string `gorm:"column:bio" json:"bio"`
	Salt         string `gorm:"column:salt" json:"-"`
	PasswordHash string `gorm:"column:password_hash" json:"-"`
	Status       string `gorm:"column:status" json:"-"`
}
