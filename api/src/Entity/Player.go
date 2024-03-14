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
	Username string `bun:",nullable" json:"username"`
	Color    string `bun:",nullable" json:"color"`
}
