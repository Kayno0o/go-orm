package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func firstOrDefault(messages []string, defaultMsg string) string {
	if len(messages) > 0 {
		return messages[0]
	}

	return defaultMsg
}

func HTTPError(c *fiber.Ctx, code int, message ...string) error {
	errorMessage := firstOrDefault(message, "Error "+string(rune(code)))
	err := c.Status(code).JSON(map[string]string{
		"error": errorMessage,
	})
	log.Println(err)
	return err
}

func HTTP400Error(c *fiber.Ctx, message ...string) error {
	return HTTPError(c, 400, firstOrDefault(message, "Bad Request"))
}

func HTTP401Error(c *fiber.Ctx, message ...string) error {
	return HTTPError(c, 401, firstOrDefault(message, "Unauthorized"))
}

func HTTP403Error(c *fiber.Ctx, message ...string) error {
	return HTTPError(c, 403, firstOrDefault(message, "Forbidden"))
}

func HTTP404Error(c *fiber.Ctx, message ...string) error {
	return HTTPError(c, 404, firstOrDefault(message, "Not Found"))
}

func HTTP500Error(c *fiber.Ctx, message ...string) error {
	return HTTPError(c, 500, firstOrDefault(message, "Internal Server Error"))
}
