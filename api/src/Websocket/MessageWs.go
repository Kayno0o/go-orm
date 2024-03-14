package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	router "go-api-test.kayn.ooo/src/Router"
)

type MessageWs struct {
	GenericWsI
}

func (ws *MessageWs) Init() {
	roomConnections := make(map[string]map[string]*Player)

	router.FiberApp.Get("/ws/message/:room", websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		uid := c.Query("uid")
		token := c.Query("token")
		user := Player{Con: c}
		user.Uid = uid
		user.Token = token

		if roomConnections[room] == nil {
			roomConnections[room] = make(map[string]*Player)
		}
		roomConnections[room][uid] = &user
		fmt.Println("new connection ", uid)

		broadcast := func(message string) {
			for userId, roomUser := range roomConnections[room] {
				if roomUser.Con == c || userId == uid {
					continue
				}

				if err := roomUser.Con.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("error writing message:", err)
					delete(roomConnections[room], userId)
				}
			}
		}

		handleClientMessage := func(message ClientMessage) {
			if message.Type == "username" {
				roomConnections[room][uid].Username = message.Content
				return
			}

			broadcast("User " + user.Username + ": " + message.Content)
		}

		defer func() {
			delete(roomConnections[room], uid)
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			message := ClientMessage{}
			if err = json.Unmarshal(msg, &message); err != nil {
				broadcast("User " + user.Username + ": " + string(msg))
				continue
			}

			handleClientMessage(message)
		}
	}))
}
