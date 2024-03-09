package middleware

import (
	"github.com/gofiber/fiber/v2"
	utils "go-api-test.kayn.ooo/src/Utils"
)

func IsLoggedIn(c *fiber.Ctx) error {
	if c.Locals("user") == nil {
		return utils.HTTP401Error(c)
	}

	return c.Next()
}
