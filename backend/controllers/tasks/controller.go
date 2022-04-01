package tasks

import (
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/middleware"
	"github.com/rafaelmf3/todo-list/models"

	"github.com/gofiber/fiber/v2"
)

type TaskService interface {
	Create(c *fiber.Ctx) error
	Read(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type taskServiceImpl struct {
	secret string
}

func NewTaskService(secret string) TaskService {
	return &taskServiceImpl{
		secret: secret,
	}
}

func (t *taskServiceImpl) Create(c *fiber.Ctx) error {
	var data models.Task

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	_, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var project models.Project
	if err := database.DB.Where("id=?", data.Project).Find(&project).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error on get Project",
		})
	}

	task := models.Task{
		Project:     project.ID,
		Description: data.Description,
	}

	database.DB.Create(&task)

	return c.JSON(task)
}

func (t *taskServiceImpl) Read(c *fiber.Ctx) error {
	_, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	id := c.Query("id")

	var task []models.Task
	if err := database.DB.Where("id=?", id).Find(&task).Error; err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "error on get a task",
		})
	}

	return c.JSON(task)
}

func (t *taskServiceImpl) Update(c *fiber.Ctx) error {
	var data models.Task

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	_, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	id := c.Query("id")

	var task models.Task
	if err := database.DB.Where("id=?", id).Find(&task).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error on get Task",
		})
	}
	if data.Description != "" {
		task.Description = data.Description
	}
	if data.Status != models.StatusNew || task.Status != models.StatusDone {
		task.Status = data.Status
	}
	database.DB.Save(&task)

	return c.JSON(task)
}

func (t *taskServiceImpl) Delete(c *fiber.Ctx) error {
	id := c.Query("id")

	_, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var task models.Task
	database.DB.Delete(&task, id)

	return c.JSON(task)
}
