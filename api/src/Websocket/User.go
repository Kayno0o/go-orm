package ws

import "github.com/gofiber/contrib/websocket"

type User struct {
	Con      *websocket.Conn `json:"-"`
	Token    string          `json:"-"`
	Id       string          `json:"id"`
	Username string          `json:"username"`
	Color    string          `json:"color"`
}

func (u *User) SendMessage(message string) error {
	return SendMessage(u.Con, message)
}
