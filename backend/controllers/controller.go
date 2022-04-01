package controllers

import (
	"github.com/rafaelmf3/todo-list/controllers/auth"
	"github.com/rafaelmf3/todo-list/controllers/lists"
	"github.com/rafaelmf3/todo-list/controllers/projects"
	"github.com/rafaelmf3/todo-list/controllers/tasks"
	"github.com/rafaelmf3/todo-list/controllers/users"
)

const SecretKey = "secret"

type Controller struct {
	User    users.UserService
	Auth    auth.AuthService
	Project projects.ProjectService
	Task    tasks.TaskService
	List    lists.ListService
}

func Start() *Controller {
	controllers := &Controller{
		User:    users.NewUserService(SecretKey),
		Auth:    auth.NewAuthService(SecretKey),
		Task:    tasks.NewTaskService(SecretKey),
		Project: projects.NewProjectService(SecretKey),
		List:    lists.NewListService(SecretKey),
	}
	return controllers
}
