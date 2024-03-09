package entity

import (
	"github.com/uptrace/bun"
	"go-api-test.kayn.ooo/src/Entity/trait"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`
	trait.Identifier
	trait.Timestampable

	Username string   `bun:",notnull" json:"username"`
	Email    string   `bun:",notnull,unique" json:"email"`
	Password string   `bun:",notnull" json:"password"`
	Roles    []string `bun:",array" json:"roles"`
}

type UserContext struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}
