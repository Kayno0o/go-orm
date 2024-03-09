package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func firstOrDefault(messages []string, defaultMsg string) string {
	if len(messages) > 0 {
		return messages[0]
	}

	return defaultMsg
}

func HTTPError(c *fiber.Ctx, code int, message ...string) error {
	err := c.Status(code).JSON(map[string]string{
		"Error": firstOrDefault(message, "Error "+string(rune(code))),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
