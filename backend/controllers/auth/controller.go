package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaelmf3/todo-list-app/backend/database"
	"github.com/rafaelmf3/todo-list-app/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authServiceImpl struct {
	secret string
}

func NewAuthService(secret string) AuthService {
	return &authServiceImpl{
		secret: secret,
	}
}

func (a *authServiceImpl) Login(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data.Email).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"massage": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(a.secret))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true, //
	}

	c.Cookie(&cookie)

	return c.JSON(&user)
}

func (a *authServiceImpl) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
