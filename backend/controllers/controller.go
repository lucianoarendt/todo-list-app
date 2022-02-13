package controllers

import (
	"github.com/rafaelmf3/todo-list-app/backend/controllers/auth"
	"github.com/rafaelmf3/todo-list-app/backend/controllers/projects"
	"github.com/rafaelmf3/todo-list-app/backend/controllers/tasks"
	"github.com/rafaelmf3/todo-list-app/backend/controllers/users"
)

const SecretKey = "secret"

type Controller struct {
	User    users.UserService
	Auth    auth.AuthService
	Project projects.ProjectService
	Task    tasks.TaskService
}

func Start() *Controller {
	controllers := &Controller{
		User:    users.NewUserService(SecretKey),
		Auth:    auth.NewAuthService(SecretKey),
		Task:    tasks.NewTaskService(SecretKey),
		Project: projects.NewProjectService(SecretKey),
	}
	return controllers
}
