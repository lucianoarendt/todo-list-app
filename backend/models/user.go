package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" gorm:"default:uuid_generate_v3()"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password []byte    `json:"-"`
}
