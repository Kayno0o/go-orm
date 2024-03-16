package ws

import utils "go-api-test.kayn.ooo/src/Utils"

type RoomI[T any] interface {
	GetUsers() map[string]*Player
	GetPublicUsers() []*Player
	AddUser(*Player)
	UpdateUser(*Player)
	Quit(*Player)
	SetAuthor(*Player)
	IsAuthor(*Player) bool
	GetUserIndex(*Player) int

	Init()

	AddMessage(*Player, string)
	HandleResponse(*Player, ClientMessage)

	Update(any)

	GetData() T
}

type Room[T any] struct {
	RoomI[T]    `json:"-"`
	Users       map[string]*Player `json:"-"`
	PublicUsers []*Player          `json:"users"`
	Messages    []Message          `json:"messages"`
	Data        T                  `json:"data"`
	Author      *Player            `json:"-"`
	Add         chan any           `json:"-"`
}

func (r *Room[T]) Init() {
	r.Users = map[string]*Player{}
	r.Messages = make([]Message, 0)
	r.Add = make(chan any)
	go Updater(r, r.Add)
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

func (r *Room[T]) AddMessage(u *Player, content string) {
	id, err := utils.RandomString(10)
	if err != nil {
		return
	}

	message := Message{u.Username, content, u.Color, id}

	r.Messages = append(r.Messages, message)

	maxMsg := 50
	if len(r.Messages) > maxMsg {
		r.Messages = r.Messages[len(r.Messages)-maxMsg:]
	}

	go Updater(r, r.Add)
	r.Add <- Update{"push", "messages", message}
}

func (r *Room[T]) Update(u any) {
	r.Add <- u
}

func (r *Room[T]) SetAuthor(u *Player) {
	r.Author = u
}

func (r *Room[T]) IsAuthor(u *Player) bool {
	return r.Author.Token == u.Token
}

func (r *Room[T]) UpdateUser(*Player) {}
