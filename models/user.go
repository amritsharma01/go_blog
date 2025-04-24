package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"` // unique constraint at DB level
	Password string `json:"password"`
}
