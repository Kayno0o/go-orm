package repository

import entity "go-api-test.kayn.ooo/src/Entity"

type UserRepositoryStruct struct {
	*GenericRepositoryStruct[entity.User]
}
