package ws

type GenericWsI interface {
	Init()
}

var (
	Users = make(map[string]*User)
)
