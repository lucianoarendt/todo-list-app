package main

import (
	"fmt"

	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var PodHash string

func main() {

	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if err := app.Listen(":8001"); err != nil {
		fmt.Println(err)
	}
}
