package entity

import (
	"github.com/uptrace/bun"
	"go-api-test.kayn.ooo/src/Entity/trait"
)

type Todolist struct {
	bun.BaseModel `bun:"table:todolist,alias:t"`
	trait.Identifier
	trait.Timestampable

	Checked bool   `bun:",notnull" json:"checked"`
	Content string `bun:",notnull" json:"content"`

	AuthorId uint `bun:",notnull" json:"author_id"`
}

type TodolistContext struct {
	ID      *uint  `json:"id"`
	Checked bool   `json:"checked"`
	Content string `json:"content"`
}
