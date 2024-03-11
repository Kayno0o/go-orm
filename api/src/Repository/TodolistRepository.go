package repository

import (
	entity "go-api-test.kayn.ooo/src/Entity"
)

type TodolistRepositoryStruct struct {
	*GenericRepositoryStruct[entity.Todolist]
}
