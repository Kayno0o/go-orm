package repository

var (
	UserRepository = &UserRepositoryStruct{}
)

type UserRepositoryStruct struct {
	GenericRepositoryStruct
}
