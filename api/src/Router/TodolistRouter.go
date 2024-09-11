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

	adminParams := Params{VerifyOwner: false, AllowPagination: true}
	adminRouter := api.Group("/admin", middleware.IsGranted([]string{"ROLE_ADMIN"}))
	adminRouter.Post("todolist", Post[entity.Todolist, entity.TodolistEditContext, entity.TodolistContext](createHandler))
	adminRouter.Put("todolist/:id", Put[entity.Todolist, entity.Todolist, entity.TodolistContext](adminParams))
	adminRouter.Delete("todolist/:id", Delete[entity.Todolist](adminParams))
	adminRouter.Get("todolist/:id", GetOne[entity.Todolist, entity.TodolistContext](adminParams))
	adminRouter.Get("todolists", GetAll[entity.Todolist, entity.TodolistContext](adminParams))
	adminRouter.Get("todolists/count", CountAll[entity.Todolist](adminParams))

	userParams := Params{VerifyOwner: true, AllowPagination: true}
	userRouter := api.Group("", middleware.IsGranted([]string{"ROLE_USER"}))
	userRouter.Post("todolist", Post[entity.Todolist, entity.TodolistEditContext, entity.TodolistContext](createHandler))
	userRouter.Put("todolist/:id", Put[entity.Todolist, entity.TodolistEditContext, entity.TodolistContext](userParams))
	userRouter.Delete("todolist/:id", Delete[entity.Todolist](userParams))
	userRouter.Get("todolist/:id", GetOne[entity.Todolist, entity.TodolistContext](userParams))
	userRouter.Get("todolists", GetAll[entity.Todolist, entity.TodolistContext](userParams))
	userRouter.Get("todolists/count", CountAll[entity.Todolist](userParams))
}
