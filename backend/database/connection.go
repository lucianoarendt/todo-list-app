package database

import (
	"fmt"

	"github.com/rafaelmf3/todo-list-app/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	conn, err := gorm.Open(mysql.Open("root:root@/db"), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Could not connect to the database, Error: %v", err))
	}

	DB = conn

	conn.AutoMigrate(&models.User{})
}
