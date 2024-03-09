package router

import (
	"github.com/gofiber/fiber/v2"
	pool "github.com/valyala/bytebufferpool"
	entity "go-api-test.kayn.ooo/src/Entity"
	repository "go-api-test.kayn.ooo/src/Repository"
	routerstruct "go-api-test.kayn.ooo/src/Router/Struct"
	"go-api-test.kayn.ooo/src/jade"
)

type ClientRouter struct {
	GenericRouterInterface
}

func getRoute(c *fiber.Ctx) routerstruct.Route {
	userInterface := c.Locals("user")
	var user *entity.User

	if userInterface != nil {
		user = userInterface.(*entity.User)
	}

	return routerstruct.Route{
		Params: c.AllParams(),
		Query:  c.Queries(),
		User:   user,
	}
}

func (cr *ClientRouter) RegisterRoutes(r fiber.Router) {
	r.Static("/public", "./public")

	r.Get(
		"/:id?",
		func(c *fiber.Ctx) error {
			bufferPool := &pool.ByteBuffer{}

			route := getRoute(c)
			route.Title = "Custom Title from Golang"

			var todolists []entity.Todolist

			if route.User != nil {
				todolists = repository.TodolistRepository.FindAllByUser(route.User)
			}

			jade.Index(todolists, route, bufferPool)
			return ServeHTML(c, bufferPool)
		},
	)
}
