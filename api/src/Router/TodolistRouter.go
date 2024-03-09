package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	group "go-api-test.kayn.ooo/src/Group"
	repository "go-api-test.kayn.ooo/src/Repository"
	utils "go-api-test.kayn.ooo/src/Utils"
	"strconv"
)

type TodolistRouter struct {
	GenericRouterInterface
}

func (tr *TodolistRouter) RegisterRoutes(r fiber.Router) {
	api := r.Group("/api")

	// ADMIN
	admin := group.IsGranted(api, []string{"ROLE_ADMIN"})
	admin.Get(
		"/todolists",
		FindAll(
			repository.TodolistRepository,
			&[]entity.Todolist{},
			&[]entity.TodolistContext{},
		),
	)

	// USER
	group.IsGranted(
		api, []string{"ROLE_USER"},
	).Post(
		"/todolist",
		tr.Post,
	).Put(
		"/todolist/:id?",
		tr.Put,
	)
}

func (tr *TodolistRouter) Put(c *fiber.Ctx) error {
	var input entity.TodolistContext
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(400)
	}

	if c.Params("id") == "" {
		return c.Status(404).SendString("Entity not found")
	}

	user := c.Locals("user").(*entity.User)
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)
	if err != nil {
		return c.Status(404).SendString("Entity not found")
	}

	todolist := &entity.Todolist{}
	err = repository.TodolistRepository.FindOneById(todolist, uint(id))
	if todolist == nil || err != nil {
		fmt.Println("error in todolist put:", err)
		return c.Status(404).SendString("Entity not found")
	}

	if todolist.AuthorId != user.ID {
		return c.Status(401).SendString("Unauthorized - Put todolist")
	}

	todolist.Checked = input.Checked

	_, err = repository.TodolistRepository.Update(todolist)
	if err != nil {
		fmt.Println(err)
	}

	output := &entity.TodolistContext{}
	utils.ApplyContext(todolist, output)

	return c.JSON(output)
}

func (tr *TodolistRouter) Post(c *fiber.Ctx) error {
	var input entity.TodolistContext
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(400)
	}

	user := c.Locals("user").(*entity.User)

	todolist := &entity.Todolist{
		AuthorId: user.ID,
		Checked:  false,
		Content:  input.Content,
	}

	_, _ = repository.TodolistRepository.Create(todolist)

	output := &entity.TodolistContext{}
	utils.ApplyContext(todolist, output)

	return c.JSON(output)
}
