package users

import (
	"gorm.io/gorm"
)

/**
Data to be modeled into the Database
*/

//TODO - Change from User to User
type User struct {
	gorm.Model
	Username     string `gorm:"column:username;unique_index"`
	Email        string `gorm:"column:email;unique_index"`
	DisplayName  string `gorm:"column:display_name"`
	Bio          string `gorm:"column:bio"`
	Salt         string `gorm:"column:salt"`
	PasswordHash string `gorm:"column:password_hash"`
	Status       string `gorm:"column:status"`
}
