package ws

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	entity "go-api-test.kayn.ooo/src/Entity"
	repository "go-api-test.kayn.ooo/src/Repository"
	router "go-api-test.kayn.ooo/src/Router"
	utils "go-api-test.kayn.ooo/src/Utils"
)

func SendMessage(c *websocket.Conn, message string) error {
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

type RoomI[T any] interface {
	GetUsers() map[string]*Player
	GetPublicUsers() []*Player
	GetUserIndex(*Player) int
	GetData() T
	Quit(*Player)
	HandleMessage(*Player, ClientMessage)
	AddUser(*Player)
	SendMessage(any)
	UpdateUser(*Player)
}

type Room[T any] struct {
	RoomI[T]    `json:"-"`
	Users       map[string]*Player `json:"-"`
	PublicUsers []*Player          `json:"users"`
	Data        T                  `json:"data"`
}

func (r *Room[T]) GetUsers() map[string]*Player {
	return r.Users
}

func (r *Room[T]) GetPublicUsers() []*Player {
	return r.PublicUsers
}

func (r *Room[T]) UpdatePublicUsers() {
	r.PublicUsers = utils.MapToArray(r.Users)
}

func (r *Room[T]) SendState() {
	r.SendUpdate("*", r)
}

func (r *Room[T]) SendUpdate(path string, object any) {
	r.SendMessage(Update{"update", path, object})
}

func (r *Room[T]) SendMessage(object any) {
	for _, u := range r.Users {
		_ = u.SendMessage(object)
	}
}

func (r *Room[T]) GetData() T {
	return r.Data
}

func (r *Room[T]) AddUser(u *Player) {
	r.Users[u.Token] = u
	r.UpdatePublicUsers()
}

func (r *Room[T]) GetUserIndex(u *Player) int {
	users := r.GetPublicUsers()
	for i, user := range users {
		if user.Token == u.Token {
			return i
		}
	}
	return -1
}

func (r *Room[T]) UpdateUser(*Player, string, any) {}

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

		if u.Username == "" {
			u.SendMessage(Update{"request", "username", nil})
		}

		var room R

		if rooms[id] == nil {
			room = init(&u)

			rooms[id] = &room
			log.Println("New:Room:" + id)
		} else {
			room = *rooms[id]
		}
		room.AddUser(&u)

		u.SendMessage(Update{"update", "*", room})
		u.SendMessage(Update{"update", "user", u})
		room.SendMessage(Update{"update", "users." + strconv.Itoa(room.GetUserIndex(&u)), u})

		defer func() {
			room.SendMessage(Update{"delete", "users." + strconv.Itoa(room.GetUserIndex(&u)), u})
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

				room.SendMessage(Update{"update", "users." + strconv.Itoa(room.GetUserIndex(&u)), u})
				u.SendMessage(Update{"update", "user", u})
				repository.Update(&u.Player, u.ID)

				room.UpdateUser(&u)
			}

			room.HandleMessage(&u, message)
		}
	}))
}
