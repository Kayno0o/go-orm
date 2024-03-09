package middleware

import (
	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	utils "go-api-test.kayn.ooo/src/Utils"
)

var RoleHierarchy = map[string][]string{
	"ROLE_USER":        {},
	"ROLE_ADMIN":       {"ROLE_USER"},
	"ROLE_SUPER_ADMIN": {"ROLE_ADMIN"},
}

func IsGranted(roles []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return utils.HTTP401Error(c)
		}

		user := localUser.(*entity.User)
		for _, role := range roles {
			if user.HasRole(role) {
				return c.Next()
			}

			for key, childRoles := range RoleHierarchy {
				if utils.Includes(childRoles, role) && user.HasRole(key) {
					return c.Next()
				}
			}
		}

		return utils.HTTP403Error(c)
	}
}
