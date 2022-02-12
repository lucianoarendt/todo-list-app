package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"default:uuid_generate_v3()"`
	Name     string
	Email    string `gorm:"unique"`
	Password []byte
}
