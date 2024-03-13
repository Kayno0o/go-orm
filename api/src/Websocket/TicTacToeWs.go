package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	router "go-api-test.kayn.ooo/src/Router"
	"log"
	"strconv"
)

type TicTacToeWs struct {
	GenericWsI
}

type messageWs struct {
	User    User   `json:"user"`
	Message string `json:"message"`
}

type messageClient struct {
	MessageType    string `json:"type"`
	MessageContent string `json:"content"`
}

func SendMessage(c *websocket.Conn, message string) error {
	return c.WriteMessage(websocket.TextMessage, []byte(message))
}

func (ws *TicTacToeWs) Init() {
	rooms := make(map[string]*TicTacToe)

	router.FiberApp.Get("/ws/tictactoe/:room", websocket.New(func(c *websocket.Conn) {
		id := c.Params("room")
		token := c.Query("token")
		var u User
		if Users[token] != nil {
			u = *Users[token]
		} else {
			u = User{
				Token: token,
				Id:    c.Query("uid"),
				Con:   c,
			}
			Users[token] = &u
			fmt.Println("New: User:" + u.Token)
		}

		var room *TicTacToe

		if rooms[id] == nil {
			room = &TicTacToe{
				Spectators: make(map[string]*User),
			}
			room.Init()

			rooms[id] = room
			fmt.Println("New: Room:" + id)
		} else {
			room = rooms[id]
		}
		rooms[id].Spectators[u.Token] = &u
		fmt.Println("New: Player:" + u.Token)

		err := u.SendMessage(room.JSON())
		if err != nil {
			return
		}

		quit := func() {
			if room.P1 != nil && room.P1.Token == u.Token {
				room.P1 = nil
				room.SendState()
			}
			if room.P2 != nil && room.P2.Token == u.Token {
				room.P2 = nil
				room.SendState()
			}
		}

		handleClientMessage := func(message messageClient) {
			switch message.MessageType {
			case "username":
				room.Spectators[u.Token].Username = message.MessageContent
				fmt.Println("Update User:" + u.Token + " Username:" + message.MessageContent)
				break
			case "click":
				var pos []int
				err := json.Unmarshal([]byte(message.MessageContent), &pos)
				if err != nil {
					fmt.Println("click error", err)
					break
				}

				if room.CurrentPlayer == 1 {
					if room.P1 != nil && room.P1.Token == u.Token {
						room.Play(pos[0], pos[1])
						room.SendState()
					}
				}

				if room.CurrentPlayer == 2 {
					if room.P2 != nil && room.P2.Token == u.Token {
						room.Play(pos[0], pos[1])
						room.SendState()
					}
				}
				break
			case "join":
				nb, err := strconv.Atoi(message.MessageContent)
				if err != nil {
					fmt.Println("join error", err)
					break
				}

				if (room.P1 != nil && room.P1.Token == u.Token) || (room.P2 != nil && room.P2.Token == u.Token) {
					fmt.Println("User:" + u.Token + " already a player")
					break
				}

				if nb == 1 {
					room.P1 = &u
					fmt.Println("User:" + u.Token + " joined as p1")
					room.SendState()
				}

				if nb == 2 {
					room.P2 = &u
					fmt.Println("User:" + u.Token + " joined as p2")
					room.SendState()
				}
			case "restart":
				if room.P1.Token == u.Token || room.P2.Token == u.Token {
					room.Init()
					room.SendState()
				}
				break
			case "quit":
				quit()
				break
			}
		}

		defer func() {
			delete(room.Spectators, u.Token)
			quit()
			fmt.Println("Del: Player:" + u.Token)
			if len(room.Spectators) == 0 {
				delete(rooms, id)
				fmt.Println("Del: Room:" + id)
			}
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
				fmt.Println("Error while receiving message", err)
				continue
			}

			handleClientMessage(message)
		}
	}))
}
