package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	router "go-api-test.kayn.ooo/src/Router"
	"log"
)

type MessageWs struct {
	GenericWsI
}

func (ws *MessageWs) Init() {
	roomConnections := make(map[string]map[string]*User)

	router.FiberApp.Get("/ws/message/:room", websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		uid := c.Query("uid")
		user := User{
			Id:  uid,
			Con: c,
		}

		if roomConnections[room] == nil {
			roomConnections[room] = make(map[string]*User)
		}
		roomConnections[room][uid] = &user
		fmt.Println("new connection ", uid)

		// Function to broadcast message to other connections in the room
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

		handleClientMessage := func(message messageClient) {
			if message.MessageType == "username" {
				roomConnections[room][uid].Username = message.MessageContent
				return
			}

			broadcast("User " + user.Username + ": " + message.MessageContent)
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

			message := messageClient{}
			if err = json.Unmarshal(msg, &message); err != nil {
				broadcast("User " + user.Username + ": " + string(msg))
				continue
			}

			handleClientMessage(message)
		}
	}))
}
