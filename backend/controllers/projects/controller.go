package projects

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rafaelmf3/todo-list-app/backend/database"
	"github.com/rafaelmf3/todo-list-app/backend/middleware"
	"github.com/rafaelmf3/todo-list-app/backend/models"
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
	return nil
}
func (t *projectServiceImpl) Update(c *fiber.Ctx) error {
	return nil
}
func (t *projectServiceImpl) Delete(c *fiber.Ctx) error {
	return nil
}
