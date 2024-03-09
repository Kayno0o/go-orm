package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	middleware "go-api-test.kayn.ooo/src/Middleware"
	repository "go-api-test.kayn.ooo/src/Repository"
	utils "go-api-test.kayn.ooo/src/Utils"
)

type TodolistRouter struct {
	GenericRouterInterface
}

func (tr *TodolistRouter) RegisterRoutes(r fiber.Router) {
	api := r.Group("/api")

	// ADMIN
	api.Get(
		"/todolist",
		middleware.IsGranted([]string{"ROLE_ADMIN"}),
		FindAll(
			repository.TodolistRepository,
			&[]entity.Todolist{},
			&[]entity.TodolistContext{},
		),
	)

	// USER
	api.Group(
		"/todolist",
		middleware.IsGranted([]string{"ROLE_USER"}),
	).Post(
		"",
		tr.Post,
	).Put(
		"/:id",
		tr.Put,
	).Get(
		"/:id",
		tr.Get,
	)
}

func (tr *TodolistRouter) Get(c *fiber.Ctx) error {
	var todolist entity.Todolist
	err := repository.FindEntityByRouteParam(c, "id", &todolist)
	if err != nil {
		return err
	}

	if !utils.IsOwner(c, todolist.AuthorId) {
		return utils.HTTP401Error(c)
	}

	return utils.JsonContext(c, todolist, &entity.TodolistContext{})
}

func (tr *TodolistRouter) Put(c *fiber.Ctx) error {
	var input entity.TodolistContext
	if err := c.BodyParser(&input); err != nil {
		return utils.HTTP400Error(c)
	}

	if c.Params("id") == "" {
		return utils.HTTP404Error(c)
	}

	var todolist entity.Todolist
	err := repository.FindEntityByRouteParam(c, "id", &todolist)

	if !utils.IsOwner(c, todolist.AuthorId) {
		return utils.HTTP401Error(c)
	}

	todolist.Checked = input.Checked

	_, err = repository.TodolistRepository.Update(todolist)
	if err != nil {
		fmt.Println(err)
	}

	return utils.JsonContext(c, todolist, &entity.TodolistContext{})
}

func (tr *TodolistRouter) Post(c *fiber.Ctx) error {
	var input entity.TodolistContext
	if err := c.BodyParser(&input); err != nil {
		return utils.HTTP400Error(c)
	}

	uid := utils.GetUserId(c)
	if uid == nil {
		return utils.HTTP401Error(c)
	}

	todolist := &entity.Todolist{
		AuthorId: *uid,
		Checked:  false,
		Content:  input.Content,
	}

	_, _ = repository.TodolistRepository.Create(todolist)

	return utils.JsonContext(c, todolist, &entity.TodolistContext{})
}
