package models

import "gorm.io/gorm"

type Symbol struct {
	gorm.Model
	Symbol string `json:"symbol"`
	ListID uint   `json:"list_id" db:"lists"`
}
