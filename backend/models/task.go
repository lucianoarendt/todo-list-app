package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	// ID          uint   `json:"id" gorm:"primaryKey"`
	Project     uint   `json:"project" db:"project"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type Status int

const (
	StatusNew  Status = 0
	StatusDone Status = 1
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}
