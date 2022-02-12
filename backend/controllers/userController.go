package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rafaelmf3/todo-list-app/backend/database"
	"github.com/rafaelmf3/todo-list-app/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// Create User
func Create(c *fiber.Ctx) error {
	var data models.User
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	user := models.User{
		ID:       uuid.New(),
		Name:     data.Email,
		Email:    data.Email,
		Password: pass,
	}

	database.DB.Create(&user)

	return c.JSON(user)
}
