package entity

import (
	"github.com/uptrace/bun"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
)

type Player struct {
	bun.BaseModel `bun:"table:player,alias:p"`
	trait.IdentifiableTrait

	Token    string `bun:",unique" json:"-"`
	Uid      string `bun:"," json:"id"`
	Username string `bun:"," json:"username"`
	Color    string `bun:"," json:"color"`
}

type Guest interface {
	Player | User
}
