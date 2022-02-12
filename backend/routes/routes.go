package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaelmf3/todo-list-app/backend/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/create", controllers.Create)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
}
