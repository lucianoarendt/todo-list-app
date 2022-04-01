package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Auth(secret string, c *fiber.Ctx) (*jwt.StandardClaims, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims, nil
}
