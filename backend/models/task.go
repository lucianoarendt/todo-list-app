package models

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Project     uint   `json:"project" db:"project"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}
