package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
)

func Includes(array []string, search string) bool {
	for _, element := range array {
		if element == search {
			return true
		}
	}

	return false
}

// ApplyContext takes and input entity and output json format
func ApplyContext(input interface{}, context interface{}) {
	jsonInput, _ := json.Marshal(input)

	err := json.Unmarshal(jsonInput, context)
	if err != nil {
		fmt.Println("ApplyContext:", err)
	}
}

func JsonContext(c *fiber.Ctx, input interface{}, context interface{}) error {
	ApplyContext(input, context)
	return c.JSON(context)
}

func GetUserId(c *fiber.Ctx) *uint {
	user := c.Locals("user").(*entity.User)
	if user == nil {
		return nil
	}
	return &user.ID
}

func IsOwner(c *fiber.Ctx, id uint) bool {
	uid := GetUserId(c)
	fmt.Println(*uid, id)
	if uid == nil || id != *uid {
		return false
	}
	return true
}
