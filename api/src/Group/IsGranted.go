package group

import (
	"github.com/gofiber/fiber/v2"
	middleware "go-api-test.kayn.ooo/src/Middleware"
)

func IsGranted(r fiber.Router, roles []string) fiber.Router {
	return r.Group(
		"",
		middleware.IsGranted(roles),
	)
}
