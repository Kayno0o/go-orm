package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	entity "go-api-test.kayn.ooo/src/Entity"
	repository "go-api-test.kayn.ooo/src/Repository"
	router "go-api-test.kayn.ooo/src/Router"
	utils "go-api-test.kayn.ooo/src/Utils"
)

func WriteWsMessage(c *websocket.Conn, message string) error {
	return c.WriteMessage(websocket.TextMessage, []byte(message))
}

type GenericWsI interface {
	Init()
}

type ClientMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type PlayerUpdate struct {
	Color    string `json:"color"`
	Username string `json:"username"`
}

type Update struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Content any    `json:"content"`
}

type UserUpdate struct {
	User   *Player
	Update Update
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Color    string `json:"color"`
	Id       string `json:"id"`
}

func Updater[T any, R RoomI[T]](r R, add chan any) {
	var updates []any

	userUpdates := make(map[*Player][]Update)

	// continuously receive updates channel
	go func() {
		defer close(add)

		for update := range add {
			switch u := update.(type) {
			case UserUpdate:
				userUpdates[u.User] = append(userUpdates[u.User], u.Update)
				break
			default:
				updates = append(updates, update)
				break
			}
		}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// on timer deliver -> send all updates in groups
	for {
		select {
		case <-ticker.C:
			if len(updates) > 0 {
				users := r.GetUsers()
				if len(users) == 0 {
					updates = []any{}
					userUpdates = make(map[*Player][]Update)
					continue
				}

				totalUpdatesLength := len(updates)
				for _, u := range userUpdates {
					totalUpdatesLength += len(u)
				}

				fmt.Println("Send:Updates:" + strconv.Itoa(totalUpdatesLength))

				for _, u := range users {
					localUpdates := make([]any, 0)

					for _, update := range userUpdates[u] {
						localUpdates = append(localUpdates, update)
					}

					for _, update := range updates {
						localUpdates = append(localUpdates, update)
					}

					if len(localUpdates) == 1 {
						_ = u.Write(localUpdates[0])
					} else {
						_ = u.Write(localUpdates)
					}
				}

				updates = []any{}
				userUpdates = make(map[*Player][]Update)
			}
		}
	}
}

func Handle[T any, R RoomI[T]](name string, init func(u *Player) R) {
	rooms := make(map[string]*R)

	router.FiberApp.Get("/ws/"+name+"/:room", websocket.New(func(c *websocket.Conn) {
		id := c.Params("room")
		token := c.Query("token")

		var u Player
		player, err := repository.FindOneBy[entity.Player](map[string]interface{}{
			"token": token,
		})

		if err != nil {
			u = Player{Con: c}
			u.Token = token
			u.Uid = c.Query("uid")

			repository.Create(&u.Player)
			log.Println("New:User:" + u.Token)
		} else {
			u = Player{
				Player: player,
				Con:    c,
			}
			log.Println("Con:User:" + utils.Stringify(u))
		}

		var room R

		if rooms[id] == nil {
			room = init(&u)
			room.Init()
			rooms[id] = &room
			room.AddUser(&u)
			room.SetAuthor(&u)
			log.Println("New:Room:" + id)
		} else {
			room = *rooms[id]
		}
		room.AddUser(&u)

		if u.Username == "" {
			room.Update(UserUpdate{&u, Update{"request", "username", nil}})
		}

		room.Update(UserUpdate{&u, Update{"update", "*", room}})
		room.Update(UserUpdate{&u, Update{"update", "user", u}})
		room.Update(Update{"update", "users." + strconv.Itoa(room.GetUserIndex(&u)), u})

		if room.IsAuthor(&u) {
			room.Update(UserUpdate{&u, Update{"update", "isAuthor", true}})
		}

		defer func() {
			room.Update(Update{"delete", "users." + strconv.Itoa(room.GetUserIndex(&u)), u})
			delete(room.GetUsers(), u.Token)
			room.Quit(&u)
			log.Println("Del:Player:" + u.Token)

			if len(room.GetUsers()) == 0 {
				delete(rooms, id)
				log.Println("Del:Room:" + id)
			}
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				break
			}
			log.Printf("Rec:Msg: %s", msg)

			message := ClientMessage{}
			if err = json.Unmarshal(msg, &message); err != nil {
				log.Println("Err:Msg:", err)
				continue
			}

			if message.Type == "user" {
				var input PlayerUpdate
				err := json.Unmarshal([]byte(message.Content), &input)
				if err != nil {
					continue
				}

				u.Username = input.Username
				color, err := utils.GetHexColor(input.Color)
				if err == nil {
					u.Color = color
				}

				room.Update(Update{"update", "users." + strconv.Itoa(room.GetUserIndex(&u)), u})
				room.Update(UserUpdate{&u, Update{"update", "user", u}})
				repository.Update(&u.Player, u.ID)

				room.UpdateUser(&u)
			}

			if message.Type == "message" {
				room.AddMessage(&u, message.Content)
				continue
			}

			room.HandleResponse(&u, message)
		}
	}))
}
