package router

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	utils "go-api-test.kayn.ooo/src/Utils"
)

type TodolistRouter struct {
	GenericRouterI
}

func (tr TodolistRouter) RegisterRoutes(r fiber.Router) {
	api := r.Group("/api")

	RegisterCrud[entity.Todolist, entity.TodolistContext, entity.TodolistEditContext](
		api,
		"todolist",
		CrudParams{},
		func(c *fiber.Ctx, context entity.TodolistEditContext) (entity.Todolist, error) {
			todo := entity.Todolist{
				Checked: false,
				Content: context.Content,
			}

			localUser := utils.GetUser(c)
			if any(localUser) == nil {
				return todo, errors.New("unauthorized")
			}

			todo.OwnerId = localUser.ID

			return todo, nil
		},
	)
}
