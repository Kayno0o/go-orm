package entity

import (
	"github.com/uptrace/bun"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
)

type Todolist struct {
	bun.BaseModel `bun:"table:todolist,alias:t"`
	trait.IdentifiableTrait
	trait.TimestampableTrait
	OwnerableTrait

	Checked bool   `bun:",notnull"`
	Content string `bun:",notnull"`
}

type TodolistContext struct {
	ID      *uint  `json:"id"`
	Checked bool   `json:"checked"`
	Content string `json:"content"`
}

type TodolistEditContext struct {
	Checked bool   `json:"checked"`
	Content string `json:"content"`
}
