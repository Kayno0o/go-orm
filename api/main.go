package main

import (
	"fmt"
	"os"

	entity "go-api-test.kayn.ooo/src/Entity"
	fixture "go-api-test.kayn.ooo/src/Fixture"
	repository "go-api-test.kayn.ooo/src/Repository"
	router "go-api-test.kayn.ooo/src/Router"
)

func main() {
	rep := repository.GenericRepository{}
	rep.Init([]interface{}{
		&entity.User{},
	})

	repository.DB.ResetModel(repository.Ctx, &entity.User{})

	if os.Getenv("ENV") == "dev" {
		count, err := repository.UserRepository.CountAll(&entity.User{})
		if err == nil && count == 0 {
			users := fixture.GenerateUsersFromJson()
			fmt.Println("Loaded", len(users), "user(s) from users.json")
		}
	}

	router.Init([]router.GenericRouterInterface{
		&router.UserRouter{},
	})
}
