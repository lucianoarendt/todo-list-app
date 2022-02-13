package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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
