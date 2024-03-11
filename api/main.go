package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	middleware "go-api-test.kayn.ooo/src/Middleware"
	"log"
	"os"

	entity "go-api-test.kayn.ooo/src/Entity"
	fixture "go-api-test.kayn.ooo/src/Fixture"
	repository "go-api-test.kayn.ooo/src/Repository"
	router "go-api-test.kayn.ooo/src/Router"
)

func main() {
	repository.Init([]interface{}{
		&entity.User{},
		&entity.Todolist{},
	})

	if os.Getenv("ENV") == "dev" {
		count, err := repository.CountAll[entity.User]()
		if err == nil && count == 0 {
			users := fixture.GenerateUsersFromJson()
			fmt.Println("Loaded", len(users), "user(s) from users.json")
		}
	}

	router.FiberApp = fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	api := router.FiberApp.Group("/", middleware.Authenticate)

	routers := []router.GenericRouterI{
		&router.UserRouter{},
		&router.TodolistRouter{},
	}

	for i := range routers {
		routers[i].RegisterRoutes(api)
	}

	log.Fatal(router.FiberApp.Listen(":3000"))
}
