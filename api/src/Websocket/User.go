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

func (u *Player) Write(message any) error {
	return WriteWsMessage(u.Con, utils.Stringify(message))
}
