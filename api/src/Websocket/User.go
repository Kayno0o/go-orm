package ws

import (
	"github.com/gofiber/contrib/websocket"
	entity "go-api-test.kayn.ooo/src/Entity"
	utils "go-api-test.kayn.ooo/src/Utils"
)

type Player struct {
	entity.Player
	Con *websocket.Conn `bun:"-" json:"-"`
}

func (u *Player) SendMessage(message any) error {
	return SendMessage(u.Con, utils.Stringify(message))
}
