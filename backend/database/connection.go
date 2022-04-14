package database

import (
	"fmt"

	"github.com/rafaelmf3/todo-list/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connectionStr := fmt.Sprintf("root:root@tcp(%s:3306)/db?parseTime=true", os.Getenv("DB_ADDRESS"))
	conn, err := gorm.Open(mysql.Open(connectionStr), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Could not connect to the database, Error: %v", err))
	}

	DB = conn

	err = conn.AutoMigrate(&models.User{}, &models.Project{}, &models.Task{}, &models.List{}, &models.Symbol{})

	if err != nil {
		panic(fmt.Sprintf("Could not migrate to the database, Error: %v", err))
	}
}
