package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaelmf3/todo-list-app/backend/database"
	"github.com/rafaelmf3/todo-list-app/backend/routes"
)

func main() {
	database.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":8000")
}
