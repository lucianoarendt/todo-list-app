package models

import "gorm.io/gorm"

type List struct {
	gorm.Model
	// ID     uint   `json:"id" gorm:"primaryKey"`
	UserID    uint     `json:"user_id"`
	IsDefault bool     `json:"default"`
	Name      string   `json:"name"`
	Symbols   []Symbol `json:"symbols" db:"symbols" gorm:"foreignKey:ListID"`
}

func (l *List) Contains(s Symbol) bool {
	for _, e := range l.Symbols {
		if e == s {
			return true
		}
	}
	return false
}
