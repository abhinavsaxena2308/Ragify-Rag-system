package models

import (
	"time"

	_ "gorm.io/gorm"
)

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Password  string     `json:"password" gorm:"not null"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Documents []Document `json:"documents" gorm:"foreignKey:UserID"`
}
