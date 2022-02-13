package models

type Project struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" db:"title"`
	Owner  uint   `json:"owner"`
	Status int    `json:"status"`
	Tasks  []Task `json:"tasks" gorm:"foreignKey:Project"`
}
