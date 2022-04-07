package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("teste", "Testando", 0).Err()

	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if err := app.Listen(":8001"); err != nil {
		fmt.Println(err)
	}
}
