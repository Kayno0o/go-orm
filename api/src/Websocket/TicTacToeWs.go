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
			Room: Room[TicTacToe]{},
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
		r.Add <- Update{"update", "data.p1", nil}
	}
	if r.Data.P2 != nil && r.Data.P2.Token == u.Token {
		r.Data.P2 = nil
		r.Add <- Update{"update", "data.p2", nil}
	}
}

func (r *TicTacToeRoom) UpdateUser(u *Player) {
	if r.Data.P1 != nil && r.Data.P1.Uid == u.Uid {
		r.Add <- Update{"update", "data.p1", u}
	}

	if r.Data.P2 != nil && r.Data.P2.Uid == u.Uid {
		r.Add <- Update{"update", "data.p2", u}
	}
}

type configUpdate struct {
	Size   int8 `json:"size"`
	Length int8 `json:"length"`
}

func (r *TicTacToeRoom) HandleResponse(u *Player, message ClientMessage) {
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
			r.Add <- Update{"update", "data.board." + strconv.Itoa(pos[1]) + "." + strconv.Itoa(pos[0]), played}
		}

		if r.Data.State != gameState {
			r.Add <- Update{"update", "data.state", r.Data.State}
		}

		if r.Data.CurrentPlayer != curr {
			r.Add <- Update{"update", "data.currentPlayer", r.Data.CurrentPlayer}
		}

		break
	case "config":
		if u.Token != r.Author.Token {
			log.Println("Err:NotAuthor:" + u.Token)
			break
		}

		var config configUpdate
		err := json.Unmarshal([]byte(message.Content), &config)
		if err != nil {
			log.Println("Err:Config:", err)
		}

		if config.Size >= 2 && config.Size <= 9 {
			r.Data.Size = config.Size
		}

		if config.Length >= 2 && config.Length <= 5 {
			r.Data.Length = config.Length
		}

		if r.Data.Size < r.Data.Length {
			r.Data.Size = r.Data.Length
		}

		r.Data.Init()

		r.Add <- Update{"update", "data.board", r.Data.Board}
		r.Add <- Update{"update", "data.state", r.Data.State}
		r.Add <- Update{"update", "data.currentPlayer", r.Data.CurrentPlayer}
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

		if nb == 1 && r.Data.P1 == nil {
			r.Data.P1 = u
			log.Println("User:" + u.Token + " joined as p1")
			r.Add <- Update{"update", "data.p1", u}
		}

		if nb == 2 && r.Data.P2 == nil {
			r.Data.P2 = u
			log.Println("User:" + u.Token + " joined as p2")
			r.Add <- Update{"update", "data.p2", u}
		}
	case "restart":
		if (r.Data.P1.Token == u.Token || r.Data.P2.Token == u.Token) && (r.Data.State == "pending" || r.Data.State == "tie") {
			r.Data.Init()
			r.Add <- Update{"update", "data", r.Data}
		}
		break
	case "quit":
		r.Quit(u)
		break
	}
}
