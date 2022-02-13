package database

import (
	"fmt"

	"github.com/rafaelmf3/todo-list-app/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	conn, err := gorm.Open(mysql.Open("root:root@/db?parseTime=true"), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Could not connect to the database, Error: %v", err))
	}

	DB = conn

	err = conn.AutoMigrate(&models.User{}, &models.Project{}, &models.Task{})

	if err != nil {
		panic(fmt.Sprintf("Could not migrate to the database, Error: %v", err))
	}
}
