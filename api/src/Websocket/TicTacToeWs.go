package ws

import (
	"encoding/json"
	"log"
	"strconv"
)

type TicTacToeWs struct {
	GenericWsI
}

func (ws *TicTacToeWs) Init() {
	Handle("tictactoe", func(u *Player) *TicTacToeRoom {
		room := &TicTacToeRoom{
			Room: Room[TicTacToe]{
				Users: map[string]*Player{
					u.Token: u,
				},
			},
		}
		room.Data.Init()
		return room
	})
}

type TicTacToeRoom struct {
	Room[TicTacToe]
}

func (r *TicTacToeRoom) Quit(u *Player) {
	if r.Data.P1 != nil && r.Data.P1.Token == u.Token {
		r.Data.P1 = nil
		r.SendUpdate("data.p1", nil)
	}
	if r.Data.P2 != nil && r.Data.P2.Token == u.Token {
		r.Data.P2 = nil
		r.SendUpdate("data.p2", nil)
	}
}

func (r *TicTacToeRoom) UpdateUser(u *Player, field string, data any) {
	if r.Data.P1 != nil && r.Data.P1.Uid == u.Uid {
		r.SendUpdate("data.p1."+field, data)
	}

	if r.Data.P2 != nil && r.Data.P2.Uid == u.Uid {
		r.SendUpdate("data.p2."+field, data)
	}
}

func (r *TicTacToeRoom) HandleMessage(u *Player, message ClientMessage) {
	switch message.Type {
	case "click":
		var pos []int
		err := json.Unmarshal([]byte(message.Content), &pos)
		if err != nil {
			log.Println("Err:Click:", err)
			break
		}

		gameState := r.Data.State
		curr := r.Data.CurrentPlayer

		played := r.Data.Play(u, pos[0], pos[1])
		if played != 0 {
			r.SendUpdate("data.board."+strconv.Itoa(pos[1])+"."+strconv.Itoa(pos[0]), played)
		}

		if r.Data.State != gameState {
			r.SendUpdate("data.state", r.Data.State)
		}

		if r.Data.CurrentPlayer != curr {
			r.SendUpdate("data.currentPlayer", r.Data.CurrentPlayer)
		}

		break
	case "join":
		nb, err := strconv.Atoi(message.Content)
		if err != nil {
			log.Println("Err:Join:", err)
			break
		}

		if (r.Data.P1 != nil && r.Data.P1.Token == u.Token) || (r.Data.P2 != nil && r.Data.P2.Token == u.Token) {
			log.Println("User:" + u.Token + " already a player")
			break
		}

		if nb == 1 {
			r.Data.P1 = u
			log.Println("User:" + u.Token + " joined as p1")
			r.SendUpdate("data.p1", u)
		}

		if nb == 2 {
			r.Data.P2 = u
			log.Println("User:" + u.Token + " joined as p2")
			r.SendUpdate("data.p2", u)
		}
	case "restart":
		if r.Data.P1.Token == u.Token || r.Data.P2.Token == u.Token {
			r.Data.Init()
			r.SendUpdate("data", r.Data)
		}
		break
	case "quit":
		r.Quit(u)
		break
	}
}
