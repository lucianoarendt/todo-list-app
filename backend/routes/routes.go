package routes

import (
	"github.com/rafaelmf3/todo-list/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	controllers := controllers.Start()
	//Auth
	app.Post("/api/login", controllers.Auth.Login)
	app.Post("/api/logout", controllers.Auth.Logout)
	//User
	app.Post("/api/user/create", controllers.User.Create)
	app.Get("/api/user/read", controllers.User.Read)
	//Project
	app.Post("/api/project/create", controllers.Project.Create)
	app.Get("/api/project/read", controllers.Project.Read)
	app.Put("/api/project/update", controllers.Project.Update)
	app.Delete("/api/project/delete", controllers.Project.Delete)
	//Tasks
	app.Post("/api/task/create", controllers.Task.Create)
	app.Get("/api/task/read", controllers.Task.Read)
	app.Put("/api/task/update", controllers.Task.Update)
	app.Delete("/api/task/delete", controllers.Task.Delete)

	//WatchLists
	app.Post("/api/list/create", controllers.List.Create)
	app.Get("/api/list/read", controllers.List.Read)
	app.Get("/api/list/readAll", controllers.List.ReadAll)
	app.Put("/api/list/update", controllers.List.Update)
	app.Delete("/api/list/delete", controllers.List.Delete)
	app.Post("/api/list/createSymbol", controllers.List.CreateSymbol)
	app.Delete("/api/list/deleteSymbol", controllers.List.DeleteSymbol)
}
