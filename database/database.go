package database

import (
	"database/sql"
	"fmt"
	"github.com/anthonyhawkins/savorbook/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func Init() *gorm.DB {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.Get("DB_HOST"),
		config.Get("DB_USER"),
		config.Get("DB_PASSWORD"),
		config.Get("DB_NAME"),
		config.Get("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("DB Error: ", err)
	}

	DB = db
	return DB

}

func GetDB() *gorm.DB {
	return DB
}

func GetSqlDB(db *gorm.DB) *sql.DB {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Unable to close DB", err)
	}
	return sqlDB
}
