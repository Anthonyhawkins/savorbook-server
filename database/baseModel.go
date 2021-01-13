package database

import (
	"time"
)

type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"userId"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}
