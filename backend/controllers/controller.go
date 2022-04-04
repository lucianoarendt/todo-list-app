package controllers

import (
	"github.com/rafaelmf3/todo-list/controllers/auth"
	"github.com/rafaelmf3/todo-list/controllers/lists/other"
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
	List    other.ListService
}

func Start() *Controller {
	controllers := &Controller{
		User:    users.NewUserService(SecretKey),
		Auth:    auth.NewAuthService(SecretKey),
		Task:    tasks.NewTaskService(SecretKey),
		Project: projects.NewProjectService(SecretKey),
		List:    other.NewListService(SecretKey),
	}
	return controllers
}
