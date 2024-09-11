package router

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	middleware "go-api-test.kayn.ooo/src/Middleware"
	utils "go-api-test.kayn.ooo/src/Utils"
)

type TodolistRouter struct {
	GenericRouterI
}

func createHandler(c *fiber.Ctx, context entity.TodolistEditContext) (entity.Todolist, error) {
	todo := entity.Todolist{
		Checked: false,
		Content: context.Content,
	}

	localUser := utils.GetUser(c)
	if localUser == nil {
		return todo, errors.New("unauthorized")
	}

	todo.OwnerId = localUser.ID

	return todo, nil
}

func (tr TodolistRouter) RegisterRoutes(r fiber.Router) {
	api := r.Group("/api")

	adminParams := Params{VerifyOwner: false, Pagination: true}
	userParams := Params{VerifyOwner: true, Pagination: true}

	adminRouter := api.Group("/admin/todolist", middleware.IsGranted([]string{"ROLE_ADMIN"}))
	adminRouter.Post("/", Post[entity.Todolist, entity.TodolistEditContext, entity.TodolistContext](createHandler))
	adminRouter.Put("/:id", Put[entity.Todolist, entity.Todolist, entity.TodolistContext](adminParams))
	adminRouter.Delete("/:id", Delete[entity.Todolist](adminParams))
	adminRouter.Get("/:id", GetOne[entity.Todolist, entity.TodolistContext](adminParams))
	adminRouter.Get("/", GetAll[entity.Todolist, entity.TodolistContext](adminParams))
	adminRouter.Get("/count", CountAll[entity.Todolist](adminParams))

	userRouter := api.Group("todolist", middleware.IsGranted([]string{"ROLE_USER"}))
	userRouter.Post("/", Post[entity.Todolist, entity.TodolistEditContext, entity.TodolistContext](createHandler))
	userRouter.Put("/:id", Put[entity.Todolist, entity.TodolistEditContext, entity.TodolistContext](userParams))
	userRouter.Delete("/:id", Delete[entity.Todolist](userParams))
	userRouter.Get("/:id", GetOne[entity.Todolist, entity.TodolistContext](userParams))
	userRouter.Get("/", GetAll[entity.Todolist, entity.TodolistContext](userParams))
	userRouter.Get("/count", CountAll[entity.Todolist](userParams))
}
