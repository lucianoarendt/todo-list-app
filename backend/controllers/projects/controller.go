package projects

import (
	"strconv"

	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/middleware"
	"github.com/rafaelmf3/todo-list/models"

	"github.com/gofiber/fiber/v2"
)

type ProjectService interface {
	Create(c *fiber.Ctx) error
	Read(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type projectServiceImpl struct {
	secret string
}

func NewProjectService(secret string) ProjectService {
	return &projectServiceImpl{
		secret: secret,
	}
}

func (t *projectServiceImpl) Create(c *fiber.Ctx) error {
	var data models.Project

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	claims, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	id, err := strconv.Atoi(claims.Issuer)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "can't convert Owner",
		})
	}

	project := models.Project{
		Owner: uint(id),
		Title: data.Title,
	}

	database.DB.Create(&project)

	return c.JSON(project)
}

func (t *projectServiceImpl) Read(c *fiber.Ctx) error {
	claims, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var project []models.Project
	if err := database.DB.Where("projects.owner=?", claims.Issuer).Find(&project).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error on read a project",
		})
	}

	for idx, i := range project {
		if err := database.DB.Where("tasks.project=?", strconv.Itoa(int(i.ID))).Find(&project[idx].Tasks).Error; err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "error on read a project's task",
			})
		}
	}

	return c.JSON(project)
}

func (t *projectServiceImpl) Update(c *fiber.Ctx) error {
	var data models.Project

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
	database.DB.Where("id=?", data.ID).Find(&project)
	project.Title = data.Title
	database.DB.Save(&project)

	return c.JSON(project)
}

func (t *projectServiceImpl) Delete(c *fiber.Ctx) error {
	id := c.Query("id")

	_, err := middleware.Auth(t.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var project models.Project
	database.DB.Delete(&project, id)

	return c.JSON(project)
}
