package models

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID           uint   `gorm:"primary_key" `
	Username     string `gorm:"column:username;unique_index"json:"username"`
	Email        string `gorm:"column:email;unique_index" json:"email"`
	DisplayName  string `gorm:"column:display_name" json:"displayName"`
	Bio          string `gorm:"column:bio" json:"bio"`
	Salt         string `gorm:"column:salt"`
	PasswordHash string `gorm:"column:password_hash" json:"password"`
	Status       string `gorm:"column:status"`
}
