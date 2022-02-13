package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rafaelmf3/todo-list-app/backend/database"
	"github.com/rafaelmf3/todo-list-app/backend/middleware"
	"github.com/rafaelmf3/todo-list-app/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(c *fiber.Ctx) error
	Read(c *fiber.Ctx) error
}

type userServiceImpl struct {
	secret string
}

func NewUserService(secret string) UserService {
	return &userServiceImpl{
		secret: secret,
	}
}

// Create User
func (u *userServiceImpl) Create(c *fiber.Ctx) error {
	var data models.User
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	user := models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: pass,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error on create user",
		})
	}

	return c.JSON(user)
}

func (u *userServiceImpl) Read(c *fiber.Ctx) error {
	claims, err := middleware.Auth(u.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var user models.User
	if err := database.DB.Joins("JOIN projects ON projects.owner = users.id").
		Where("users.id=?", claims.Issuer).
		Find(&user).
		Select("projects.id, title, owner, projects.status").
		Find(&user.Projects); err != nil {
		fmt.Println(err)
	}

	return c.JSON(user)
}
