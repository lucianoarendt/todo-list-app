package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	// ID       uint      `json:"id" db:"id" gorm:"primaryKey"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email" gorm:"unique"`
	Password []byte    `json:"-" db:"-"`
	Projects []Project `json:"projects" db:"projects" gorm:"foreignKey:Owner"`
}
