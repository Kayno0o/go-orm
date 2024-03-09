package repository

import (
	"fmt"
	entity "go-api-test.kayn.ooo/src/Entity"
)

var (
	TodolistRepository = &TodolistRepositoryStruct{}
)

type TodolistRepositoryStruct struct {
	GenericRepository
}

func (tr TodolistRepositoryStruct) FindAllByUser(user *entity.User) []entity.Todolist {
	var todolists []entity.Todolist

	err := tr.FindAllBy(&todolists, map[string]interface{}{
		"AuthorId": user.ID,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return todolists
}
