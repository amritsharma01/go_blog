package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AuthorID    uint     `json:"author_id"`
	Author      User     `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	CategoryID  uint     `json:"category_id" gorm:"default:6"`
	Category    Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
}
